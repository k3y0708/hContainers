package main

import (
	"fmt"
	"os"

	"github.com/hContainers/hContainers/assembler"
	"github.com/hContainers/hContainers/printer"
	"github.com/hContainers/hContainers/service"
	"github.com/hContainers/hContainers/util"
)

func cliLocations(args []string) {
	if len(args) == 0 {
		cliLocationsHelp()
		os.Exit(1)
	}
	switch args[0] {
	case "list":
		printer.LocationsList(assembler.LocationsToStrings(service.GetAllLocations()))
	default:
		fmt.Println("Invalid command")
		nearestCommand := util.FindNearestCommand(args[0], []string{"list"})
		if nearestCommand != "" {
			fmt.Println("Did you mean 'hContainers locations " + nearestCommand + "'?")
		} else {
			cliLocationsHelp()
		}
	}
}
