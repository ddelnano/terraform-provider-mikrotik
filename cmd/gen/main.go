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
	"strings"

	"github.com/ddelnano/terraform-provider-mikrotik/cmd/gen/internal/codegen"
)

type (
	Configuration struct {
		SrcFile        string
		DestFile       string
		StructName     string
		IDFieldName    string
		RequiredFields map[string]bool
		OptionalFields map[string]bool
		ComputedFields map[string]bool
		OmitFields     map[string]bool
	}
)

func main() {
	if err := realMain(os.Args[1:]); err != nil {
		log.Fatalf("execution failed: %v", err)
	}

	os.Exit(0)
}

// todo:
// 		optional fields
// 		omit fields
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
		requiredFields = map[string]bool{}
		optionalFields = map[string]bool{}
		computedFields = map[string]bool{}
		omitFields     = map[string]bool{}
	)
	flag.Func("requiredFields", "A comma separated list of required fields", func(s string) error {
		for _, v := range strings.Split(s, ",") {
			requiredFields[strings.ToLower(v)] = true
		}
		return nil
	})

	flag.Func("optionalFields", "A comma separated list of optional fields", func(s string) error {
		for _, v := range strings.Split(s, ",") {
			optionalFields[strings.ToLower(v)] = true
		}
		return nil
	})

	flag.Func("computedFields", "A comma separated list of computed fields", func(s string) error {
		for _, v := range strings.Split(s, ",") {
			computedFields[strings.ToLower(v)] = true
		}
		return nil
	})
	flag.Func("omitFields", "A comma separated list of fields to omit", func(s string) error {
		for _, v := range strings.Split(s, ",") {
			omitFields[strings.ToLower(v)] = true
		}
		return nil
	})

	if err := flag.CommandLine.Parse(args); err != nil {
		return err
	}

	config := Configuration{}
	config.DestFile = *destFile
	config.SrcFile = *srcFile
	config.IDFieldName = *idField
	config.StructName = *structName
	config.IDFieldName = *idField
	config.RequiredFields = requiredFields
	config.OmitFields = omitFields
	config.OptionalFields = optionalFields
	config.ComputedFields = computedFields

	if config.SrcFile == "" {
		var err error
		config.SrcFile, err = toAbsPath(os.Getenv("GOFILE"), ".")
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
	// we delay this initialisation, because struct name might be available only after file parsing
	if *destFile == "" {
		config.DestFile, err = toAbsPath(path.Join("../mikrotik", structNameToResourceFilename(config.StructName)), "./")
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
	var isPrevLower bool
	var buf strings.Builder

	for _, r := range structName {
		if 'A' <= r && r <= 'Z' && isPrevLower {
			buf.WriteByte('_')
			buf.WriteString(strings.ToLower(string(r)))
			isPrevLower = false
			continue
		}

		isPrevLower = 'a' <= r && r <= 'z'
		buf.WriteString(strings.ToLower(string(r)))
	}

	return "resource_" + buf.String() + ".go"
}
