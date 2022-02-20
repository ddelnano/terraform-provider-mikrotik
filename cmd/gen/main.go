package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"

	"github.com/ddelnano/terraform-provider-mikrotik/cmd/gen/internal/codegen"
)

func main() {
	os.Exit(realMain(os.Args[1:]))
}

func realMain(args []string) int {
	// filename := "client/dns.go"
	filename := args[0]
	if err := processFile(filename); err != nil {
		log.Print(err)
		return 1
	}

	return 0
}

func processFile(filename string) error {
	fSet := token.NewFileSet()
	aFile, err := parser.ParseFile(fSet, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	if aFile == nil {
		return errors.New("parsing of the file failed")
	}

	// ast.Print(fSet, aFile)
	structName := "DnsRecord"
	s, err := codegen.Parse(aFile, structName)
	if err != nil {
		return err
	}

	if err := generateResource(*s, os.Stdout); err != nil {
		return err
	}

	return nil
}

func generateResource(s codegen.Struct, w io.Writer) error {
	fmt.Fprintln(w, "//"+s.Name)
	for _, v := range s.Fields {
		fmt.Fprintf(w, "//    %s\t%s\t%s\n", v.Name, v.Type, v.Tag)
	}
	fmt.Fprintf(w, "//=====================================================\n")

	buf := bytes.Buffer{}
	if err := codegen.WriteSource(&buf, s); err != nil {
		return err
	}

	var err error
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	// formatted := buf.Bytes()
	_, err = w.Write(formatted)
	if err != nil {
		return err
	}

	return nil
}
