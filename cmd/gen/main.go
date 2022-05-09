package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/ddelnano/terraform-provider-mikrotik/cmd/gen/internal/codegen"
	"github.com/ddelnano/terraform-provider-mikrotik/cmd/gen/internal/utils"
)

type (
	Configuration struct {
		SrcFile     string
		DestFile    string
		StructName  string
		IDFieldName string
	}
)

func main() {
	if err := realMain(os.Args[1:]); err != nil {
		log.Fatalf("execution failed: %v", err)
	}

	os.Exit(0)
}

// todo:
// 		omit fields by default: Id
// 		different fields for update/delete client funcs
// 			(sometimes it is Id, sometimes it's Name, etc)
func realMain(args []string) error {
	var (
		destFile       = flag.String("dest", "", "File to write result to")
		srcFile        = flag.String("src", "", "Source file to parse struct from")
		structName     = flag.String("struct", "", "Name of a struct to process")
		idField        = flag.String("idField", "Id", "Name of a struct field to use as Terraform ID of resource")
		skipFormatting = flag.Bool("skipFormatting", false, "Whether code formatting should be skipped")
	)

	if err := flag.CommandLine.Parse(args); err != nil {
		return err
	}

	config := Configuration{}
	config.DestFile = *destFile
	config.SrcFile = *srcFile
	config.IDFieldName = *idField
	config.StructName = *structName
	config.IDFieldName = *idField

	if config.SrcFile == "" {
		var err error
		config.SrcFile, err = filepath.Abs(os.Getenv("GOFILE"))
		if err != nil {
			return err
		}
	}

	if config.IDFieldName == "" {
		return errors.New("idField name must be present")
	}

	startLine := 1
	lineStr := os.Getenv("GOLINE")
	if lineStr != "" {
		lineInt, err := strconv.Atoi(lineStr)
		if err != nil {
			return fmt.Errorf("fail to parse GOLINE: %v", err.Error())
		}
		startLine = lineInt
	}
	s, err := processFile(config.SrcFile, startLine, config.StructName)
	if err != nil {
		return err
	}
	if config.StructName == "" {
		config.StructName = s.Name
	}
	// we delay this initialization, because struct name might be available only after file parsing
	if *destFile == "" {
		config.DestFile, err = filepath.Abs(path.Join("../mikrotik", "resource_"+utils.ToSnakeCase(config.StructName)) + ".go")
		if err != nil {
			return err
		}
	}
	if config.DestFile == "" {
		return errors.New("destination file must be set via flags or 'go:generate' mode must be used")
	}

	s.IDFieldName = config.IDFieldName

	var out io.Writer
	if config.DestFile == "-" {
		out = os.Stdout
	} else {
		file, err := os.OpenFile(config.DestFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		out = file
		defer func() {
			file.Close()
		}()
	}
	return generateResource(s, out, !*skipFormatting)
}

func processFile(filename string, startLine int, structName string) (*codegen.Struct, error) {
	_, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	fSet := token.NewFileSet()
	aFile, err := parser.ParseFile(fSet, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	if aFile == nil {
		return nil, errors.New("parsing of the file failed")
	}

	s, err := codegen.Parse(fSet, aFile, startLine, structName)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func generateResource(s *codegen.Struct, w io.Writer, formatCode bool) error {
	var result []byte
	var buf bytes.Buffer

	if err := codegen.WriteSource(&buf, *s); err != nil {
		return err
	}
	result = buf.Bytes()

	if formatCode {
		var err error
		result, err = format.Source(buf.Bytes())
		if err != nil {
			return err
		}
	}

	_, err := w.Write(result)
	if err != nil {
		return err
	}

	return nil
}
