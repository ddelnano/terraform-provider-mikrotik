package consoleinspected

import "strings"

var DefaultSplitStrategy = splitStrategyFunc(orderedSplit)

type splitStrategyFunc func(string) ([]string, error)

func (f splitStrategyFunc) Split(in string) ([]string, error) {
	return f(in)
}

// orderedSplit splits items definition using order of fields.
//
// Each 'name=' key starts a new item definition.
func orderedSplit(in string) ([]string, error) {
	result := []string{}

	buf := strings.Builder{}
	for _, v := range strings.Split(in, ";") {
		if strings.TrimSpace(v) == "" {
			continue
		}
		if strings.HasPrefix(v, "name=") {
			if buf.Len() > 0 {
				result = append(result, buf.String())
			}
			buf.Reset()
		}
		buf.WriteString(v)
		buf.WriteString(";")
	}
	if buf.Len() > 0 {
		result = append(result, buf.String())
	}

	return result, nil
}
