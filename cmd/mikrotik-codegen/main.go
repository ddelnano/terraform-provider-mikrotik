package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ddelnano/terraform-provider-mikrotik/cmd/mikrotik-codegen/internal/codegen"
)

type (
	Configuration struct {
		SrcFile    string
		DestFile   string
		StructName string
		FormatCode bool
	}
)

func main() {
	if err := realMain(os.Args[1:]); err != nil {
		log.Fatalf("execution failed: %v", err)
	}

	os.Exit(0)
}

// todo:
//
//	omit fields by default: Id
//	different fields for update/delete client funcs
//		(sometimes it is Id, sometimes it's Name, etc)
func realMain(args []string) error {
	config, err := parseConfig(args)
	if err != nil {
		return err
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

	s, err := codegen.ParseFile(config.SrcFile, startLine, config.StructName)
	if err != nil {
		return err
	}
	if s.TerraformIDField == "" {
		return errors.New("Terraform ID field is not set to any of struct's fields")
	}
	if s.MikrotikIDField == "" {
		return errors.New("Mikrotik ID field is not set to any of struct's fields")
	}

	if config.StructName == "" {
		config.StructName = s.Name
	}
	// if destination was not provided via args, build it using struct name
	// the initialization is delayed, because config.StructName might not be provided, so we get it after actual file parsing
	// skip this not to accidently rewrite existing file
	// if config.DestFile == "" {
	// 	config.DestFile, err = filepath.Abs(path.Join("../mikrotik", "resource_"+utils.ToSnakeCase(config.StructName)) + ".go")
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	if config.DestFile == "" {
		return errors.New("destination file must be set via flags or 'go:generate' mode must be used")
	}

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
	writeHooks := []codegen.SourceWriteHookFunc{}
	if config.FormatCode {
		writeHooks = append(writeHooks, codegen.SourceFormatHook)
	}
	return codegen.GenerateResource(s, out, writeHooks...)
}

func parseConfig(args []string) (*Configuration, error) {
	var (
		destFile   = flag.String("dest", "-", "File to write result to. Default: write to stdout.")
		srcFile    = flag.String("src", "", "Source file to parse struct from.")
		structName = flag.String("struct", "", "Name of a struct to process.")
		formatCode = flag.Bool("formatCode", true, "Whether to format resulting code.")
	)

	if err := flag.CommandLine.Parse(args); err != nil {
		return nil, err
	}

	config := Configuration{}
	config.DestFile = *destFile
	config.SrcFile = *srcFile
	config.StructName = *structName
	config.FormatCode = *formatCode

	if config.SrcFile == "" {
		var err error
		config.SrcFile, err = filepath.Abs(os.Getenv("GOFILE"))
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}
