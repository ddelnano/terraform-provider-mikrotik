package consoleinspected

const (
	// NodeTypeDir represents console menu level.
	NodeTypeDir NodeType = "dir"

	// NodeTypeCommand represents console command that can be called.
	NodeTypeCommand NodeType = "cmd"

	// NodeTypeArg represents console item that is argument to a command.
	NodeTypeArg NodeType = "arg"

	// TypeSelf is console item type for currently inspected item.
	TypeSelf Type = "self"

	// TypeChild is console item type of all items within inspected container.
	TypeChild Type = "child"
)

type (
	// NodeType is dedicated type that holds values of "node-type" field of console item.
	NodeType string

	// Type is dedicated type that holds values of "type" field of console item.
	Type string

	// Item represents inspected console items.
	Item struct {
		NodeType NodeType `mikrotik:"node-type"`
		Type     Type     `mikrotik:"type"`
		Name     string   `mikrotik:"name"`
	}

	// ConsoleItem represents inspected console item with extracted commands, arguments, etc.
	ConsoleItem struct {
		// Self holds information about current console item.
		Self Item

		// Commands holds a list of commands available for this menu level.
		// Valid only for ConsoleItem of type NodeTypeDir.
		Commands []string

		// Subcommands holds a list of commands for the nested menu level.
		// Valid only for ConsoleItem of type NodeTypeDir.
		Subcommands []string

		// Arguments holds a list of argument items for a command.
		// Valid only for ConsoleItem of type NodeItemCommand.
		Arguments []Item
	}
)

type (
	ItemsDefinitionSplitStrategy interface {
		// Split splits set of items definition represented by a single string into chunks of separate item definitions.
		Split(string) ([]string, error)
	}
)
