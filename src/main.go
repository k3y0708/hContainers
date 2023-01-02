package main

import (
	"fmt"
	"os"

	"github.com/hContainers/hContainers/util"

	"github.com/hContainers/hContainers/global"
)

func main() {
	args := os.Args[1:]

	util.CheckLength(args, 1, "No arguments provided\nTry 'hContainers help' for more information.", 1)

	envvarlessCommands := []string{"help", "version"}
	if !util.Contains(envvarlessCommands, args[0]) {
		util.GetEnv()
	}

	switch args[0] {
	case "help":
		cliHelp()
	case "runners":
		cliRunner(args[1:])
	case "locations":
		cliLocations(args[1:])
	case "containers":
		cliContainers(args[1:])
	case "version":
		fmt.Println(global.Version)

	default:
		fmt.Println("Invalid command")
		nearestCommand := util.FindNearestCommand(args[0], []string{"help", "runners", "locations", "containers", "version"})
		if nearestCommand != "" {
			fmt.Println("Did you mean 'hContainers " + nearestCommand + "'?")
		} else {
			cliHelp()
		}
	}
}
