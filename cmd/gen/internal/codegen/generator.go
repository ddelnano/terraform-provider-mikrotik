package codegen

import (
	"bytes"
	"io"
)

type (
	// SourceWriteHookFunc defines a hook func to mutate source before writing to destination
	SourceWriteHookFunc func([]byte) ([]byte, error)
)

// GenerateResource generates Terraform resource and writes it to specified output
func GenerateResource(s *Struct, w io.Writer, beforeWriteHooks ...SourceWriteHookFunc) error {
	var result []byte
	var buf bytes.Buffer
	var err error

	if err := WriteSource(&buf, *s); err != nil {
		return err
	}
	result = buf.Bytes()
	for _, h := range beforeWriteHooks {
		result, err = h(result)
		if err != nil {
			return err
		}
	}

	_, err = w.Write(result)
	if err != nil {
		return err
	}

	return nil
}
