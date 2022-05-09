package codegen

import (
	"go/format"
)

// SourceFormatHook formats code using Go's formatter
func SourceFormatHook(p []byte) ([]byte, error) {
	return format.Source(p)
}
