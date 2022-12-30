package main

import (
	"os"

	"github.com/hContainers/hContainers/service"
)

func cliLocations(args []string) {
	if len(args) == 0 {
		cliLocationsHelp()
		os.Exit(1)
	}
	switch args[0] {
	case "list":
		service.LocationList()
	default:
		cliLocationsHelp()
	}
}
