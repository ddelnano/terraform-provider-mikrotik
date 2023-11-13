package consoleinspected

import (
	"fmt"
	"strings"
)

// Parse parses definition of inspected console item and extracts items using splitStrategy.
//
// It returns console item struct with its subcommands, commands, arguments, etc.
func Parse(input string, splitStrategy ItemsDefinitionSplitStrategy) (ConsoleItem, error) {
	chunks, err := splitStrategy.Split(input)
	if err != nil {
		return ConsoleItem{}, err
	}

	var result ConsoleItem
	result.Self = Item{}

	for _, v := range chunks {
		item, err := parseItem(v)
		if err != nil {
			return ConsoleItem{}, err
		}
		if item.Type == TypeSelf {
			result.Self = item
			continue
		}
		switch t := item.NodeType; t {
		case NodeTypeDir:
			result.Subcommands = append(result.Subcommands, item.Name)
		case NodeTypeArg:
			result.Arguments = append(result.Arguments, item)
		case NodeTypeCommand:
			result.Commands = append(result.Commands, item.Name)
		default:
			return ConsoleItem{}, fmt.Errorf("unknown node type %q", t)
		}
	}

	return result, nil
}

func parseItem(input string) (Item, error) {
	result := Item{}
	for _, v := range strings.Split(input, ";") {
		if strings.TrimSpace(v) == "" {
			continue
		}
		if strings.HasPrefix(v, "name=") {
			result.Name = strings.TrimPrefix(v, "name=")
		}
		if strings.HasPrefix(v, "node-type=") {
			result.NodeType = NodeType(strings.TrimPrefix(v, "node-type="))
		}
		if strings.HasPrefix(v, "type=") {
			result.Type = Type(strings.TrimPrefix(v, "type="))
		}
	}

	return result, nil
}
