package consoleinspected

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expected      ConsoleItem
		expectedError bool
	}{
		{
			name:  "simple command",
			input: "name=add;node-type=cmd;type=self;name=comment;node-type=arg;type=child;name=copy-from;node-type=arg;type=child;",
			expected: ConsoleItem{
				Self: Item{
					Name:     "add",
					NodeType: NodeTypeCommand,
					Type:     TypeSelf,
				},
				Arguments: []Item{
					{Name: "comment", NodeType: NodeTypeArg, Type: TypeChild},
					{Name: "copy-from", NodeType: NodeTypeArg, Type: TypeChild},
				},
			},
		},
		{
			name:  "command with subcommands",
			input: "name=list;node-type=dir;type=self;name=add;node-type=cmd;type=child;name=comment;node-type=cmd;type=child;name=edit;node-type=cmd;type=child;name=export;node-type=cmd;type=child;name=find;node-type=cmd;type=child;name=get;node-type=cmd;type=child;name=member;node-type=dir;type=child;name=print;node-type=cmd;type=child;name=remove;node-type=cmd;type=child;name=reset;node-type=cmd;type=child;name=set;node-type=cmd;type=child",
			expected: ConsoleItem{
				Self: Item{
					Name:     "list",
					NodeType: NodeTypeDir,
					Type:     TypeSelf,
				},
				Commands: []string{
					"add",
					"comment",
					"edit",
					"export",
					"find",
					"get",
					"print",
					"remove",
					"reset",
					"set",
				},
				Subcommands: []string{"member"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			item, err := Parse(tc.input, DefaultSplitStrategy)
			if !assert.Equal(t, !tc.expectedError, err == nil) || tc.expectedError {
				return
			}
			assert.Equal(t, tc.expected, item)
		})
	}
}
