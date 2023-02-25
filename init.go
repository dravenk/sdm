package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

var (
	InputUp     = "up"
	InputInit   = "init"
	InputDown   = "down"
	InputRemove = "remove"
)

func defineCommand(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: InputInit, Description: "Initialize the default configuration."},
		{Text: InputUp, Description: "Create and start containers."},
		{Text: InputDown, Description: "Stop and remove containers, networks."},
		{Text: InputRemove, Description: "Remove all applications"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func initCommand() string {
	fmt.Println("Please type you command.")
	t := prompt.Input("> ", defineCommand)
	switch t {
	case InputInit:
		initConfigFile()
		initDockerComposefile()
		initSettingsfile()
		return InputInit
	case InputUp:
		return InputUp
	case InputDown:
		return InputDown
	case InputRemove:
		return InputRemove
	default:
		return ""
	}
}
