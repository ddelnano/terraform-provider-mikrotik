package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	consoleinspected "github.com/ddelnano/terraform-provider-mikrotik/client/console-inspected"
	"github.com/ddelnano/terraform-provider-mikrotik/cmd/mikrotik-codegen/internal/codegen"
)

type (
	MikrotikConfiguration struct {
		CommandBasePath       string
		ResourceName          string
		Test                  bool
		SrcFile               string
		InspectDefinitionFile string
		QueryDefinition       bool
	}

	TerraformConfiguration struct {
		SrcFile    string
		StructName string
		AccTest    bool
	}

	GeneratorFunc func(w io.Writer) error
)

func main() {
	if err := realMain(os.Args[1:]); err != nil {
		log.Fatalf("execution failed: %v", err)
	}

	os.Exit(0)
}

func usage(w io.Writer) {
	_, _ = w.Write([]byte(`
Sub commands:
  mikrotik - generate MikroTik model
  terraform - generate Terraform resource
`))
}

func realMain(args []string) error {
	if len(args) < 1 {
		usage(flag.CommandLine.Output())
		return nil
	}
	subcommand := args[0]
	args = args[1:]

	var formatCode bool
	var destFile string
	var generator func() GeneratorFunc

	switch subcommand {
	case "terraform":
		config := TerraformConfiguration{}
		fs := flag.NewFlagSet("terraform", flag.ExitOnError)
		commonFlags(fs, &destFile, &formatCode)
		fs.StringVar(&config.SrcFile, "src", "", "Source file to parse struct from.")
		fs.StringVar(&config.StructName, "struct", "", "Name of a struct to process.")
		fs.BoolVar(&config.AccTest, "accTest", false, "Generate acceptance test instead.")
		_ = fs.Parse(args)

		startLine := 1
		lineStr := os.Getenv("GOLINE")
		if lineStr != "" {
			lineInt, err := strconv.Atoi(lineStr)
			if err != nil {
				return fmt.Errorf("fail to parse GOLINE: %v", err.Error())
			}
			startLine = lineInt
		}

		s, err := parseFile(config.SrcFile, startLine, config.StructName)
		if err != nil {
			return err
		}

		// If struct name is not provider, use one found in the parsed file.
		// See `ParseFile()` for details.
		if config.StructName == "" {
			config.StructName = s.Name
		}

		if destFile == "" {
			return errors.New("destination file must be set via flags or 'go:generate' mode must be used")
		}

		generator = func(s *codegen.Struct) func() GeneratorFunc {
			return func() GeneratorFunc {
				return func(w io.Writer) error {
					return codegen.GenerateResource(s, w)
				}
			}
		}(s)

		if config.AccTest {
			generator = func(s *codegen.Struct) func() GeneratorFunc {
				return func() GeneratorFunc {
					return func(w io.Writer) error {
						return codegen.GenerateResourceTest(s, w)
					}
				}
			}(s)

		}

	case "mikrotik":
		config := MikrotikConfiguration{}
		fs := flag.NewFlagSet("mikrotik", flag.ExitOnError)
		commonFlags(fs, &destFile, &formatCode)

		fs.StringVar(&config.ResourceName, "name", "", "Name of the resource to generate main code or test-file.")
		fs.StringVar(&config.CommandBasePath, "commandBase", "", "The command base path in MikroTik.")
		fs.StringVar(&config.SrcFile, "src", "", "Source file to parse struct from. Conflicts with 'commandBase'")
		fs.BoolVar(&config.Test, "test", false, "Generate resource test-file instead.")
		fs.StringVar(&config.InspectDefinitionFile, "inspect-definition-file", "",
			"[EXPERIMENTAL] File with command definition. Conflicts with query-definition.")
		fs.BoolVar(&config.QueryDefinition, "query-definition", false,
			"[EXPERIMENTAL] Query remote MikroTik device to fetch resource fields. Conflicts with inspect-definition-file.")

    _ = fs.Parse(args)

		if config.InspectDefinitionFile != "" && config.QueryDefinition {
			return errors.New("only one of inspect-definition-file or query-definition can be used")
		}

		consoleCommandDefinition := consoleinspected.ConsoleItem{}
		if config.InspectDefinitionFile != "" {
			fileBytes, err := os.ReadFile(config.InspectDefinitionFile)
			if err != nil {
				return err
			}

			consoleCommandDefinition, err = consoleinspected.Parse(string(fileBytes), consoleinspected.DefaultSplitStrategy)
			if err != nil {
				return err
			}
		}

		if config.QueryDefinition {
			var err error
			c := client.NewClient(client.GetConfigFromEnv())
			consoleCommandDefinition, err = c.InspectConsoleCommand(config.CommandBasePath + "/add")
			if err != nil {
				return err
			}
		}

		generator = func() GeneratorFunc {
			return func(w io.Writer) error {
				return codegen.GenerateMikrotikResource(config.ResourceName, config.CommandBasePath, consoleCommandDefinition, w)
			}
		}
		if config.Test {
			if config.CommandBasePath != "" {
				return errors.New("while generating test-file, 'commandBase' flags must not be set")
			}
			if config.SrcFile == "" {
				return errors.New("in test-file generating mode, 'src' flag must point to source file with struct.")
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

			s, err := parseFile(config.SrcFile, startLine, config.ResourceName)
			if err != nil {
				return err
			}

			generator = func() GeneratorFunc {
				return func(w io.Writer) error {
					return codegen.GenerateMikrotikResourceTest(config.ResourceName, s, w)
				}
			}
		}

	default:
		return errors.New("unsupported subcommand: " + subcommand)
	}

	var out io.Writer
	if destFile == "-" {
		out = os.Stdout
	} else {
		file, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		out = file
		defer func() {
			file.Close()
		}()
	}
	writeHooks := []codegen.SourceWriteHookFunc{}
	if formatCode {
		writeHooks = append(writeHooks, codegen.SourceFormatHook)
	}

	var err error
	var buf bytes.Buffer

	if err := generator()(&buf); err != nil {
		return err
	}

	var result []byte
	result = buf.Bytes()
	for _, h := range writeHooks {
		result, err = h(result)
		if err != nil {
			return err
		}
	}

	if _, err := out.Write(result); err != nil {
		return err
	}

	return nil
}

func commonFlags(fs *flag.FlagSet, dest *string, formatCode *bool) {
	fs.StringVar(dest, "dest", "-", "File to write result to. Default: write to stdout.")
	fs.BoolVar(formatCode, "formatCode", true, "Whether to format resulting code. Useful for debugging to see raw source code right after generation.")
}

func parseFile(srcFile string, startLine int, structName string) (*codegen.Struct, error) {
	s, err := codegen.ParseFile(srcFile, startLine, structName)
	if err != nil {
		return nil, err
	}

	return s, nil
}
