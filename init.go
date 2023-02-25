package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

var (
	InputUp   = "up"
	InputInit = "init"
	InputDown = "down"
)

func defineCommand(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: InputInit, Description: "Initialize configuration by default."},
		{Text: InputUp, Description: "Excute docker-compose to create and start containers."},
		{Text: InputDown, Description: "Excute docker-compose down to stop and remove containers, networks"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func initCommand() string {
	fmt.Println("Please type you command: init, up or down.")
	t := prompt.Input("> ", defineCommand)
	// fmt.Println("You selected " + t)
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
	default:
		return ""
	}
}
