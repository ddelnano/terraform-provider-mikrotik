package client

import (
	"strings"

	consoleinspected "github.com/ddelnano/terraform-provider-mikrotik/client/console-inspected"
)

func (c Mikrotik) InspectConsoleCommand(command string) (consoleinspected.ConsoleItem, error) {
	client, err := c.getMikrotikClient()
	if err != nil {
		return consoleinspected.ConsoleItem{}, err
	}
	normalizedCommand := strings.ReplaceAll(command[1:], "/", ",")
	cmd := []string{"/console/inspect", "as-value", "=path=" + normalizedCommand, "=request=child"}
	reply, err := client.RunArgs(cmd)
	if err != nil {
		return consoleinspected.ConsoleItem{}, err
	}
	var items []consoleinspected.Item
	var result consoleinspected.ConsoleItem
	if err := Unmarshal(*reply, &items); err != nil {
		return consoleinspected.ConsoleItem{}, err
	}

	for _, v := range items {
		if v.Type == consoleinspected.TypeSelf {
			result.Self = v
			continue
		}
		switch v.NodeType {
		case consoleinspected.NodeTypeArg:
			result.Arguments = append(result.Arguments, v)
		case consoleinspected.NodeTypeCommand:
			result.Commands = append(result.Commands, v.Name)
		case consoleinspected.NodeTypeDir:
			result.Subcommands = append(result.Subcommands, v.Name)
		}
	}

	return result, nil
}
