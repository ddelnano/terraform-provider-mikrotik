package main

import (
	"bytes"
	"errors"
	"flag"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/ddelnano/terraform-provider-mikrotik/cmd/gen/internal/codegen"
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

	if config.StructName == "" {
		return errors.New("struct name must be set")
	}
	if config.SrcFile == "" {
		var err error
		config.SrcFile, err = toAbsPath(os.Getenv("GOFILE"), ".")
		if err != nil {
			return err
		}
		config.DestFile, err = toAbsPath(path.Join("../mikrotik", structNameToResourceFilename(config.StructName)), "./")
		if err != nil {
			return err
		}
	}
	if config.DestFile == "" {
		return errors.New("destination file must be set via flags or 'go:generate' mode must be used")
	}
	if config.IDFieldName == "" {
		return errors.New("idField name must be present")
	}

	s, err := processFile(config.SrcFile, config.StructName)
	if err != nil {
		return err
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

func processFile(filename, structName string) (*codegen.Struct, error) {
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

	s, err := codegen.Parse(aFile, structName)
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

func toAbsPath(filename string, workdirs ...string) (string, error) {
	if path.IsAbs(filename) {
		return filename, nil
	}

	absPath := filename
	for _, w := range workdirs {
		if len(w) > 0 {
			absPath = path.Join(w, filename)
			break
		}
	}

	return filepath.Abs(absPath)
}

func structNameToResourceFilename(structName string) string {
	return structName
}
