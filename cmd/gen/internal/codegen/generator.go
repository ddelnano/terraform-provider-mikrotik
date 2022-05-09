package codegen

import (
	"bytes"
	"go/format"
	"io"
)

// GenerateResource generates Terraform resource and writes it to specified output
func GenerateResource(s *Struct, w io.Writer, formatCode bool) error {
	var result []byte
	var buf bytes.Buffer

	if err := WriteSource(&buf, *s); err != nil {
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
