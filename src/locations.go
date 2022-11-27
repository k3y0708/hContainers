package main

import (
	"k3y0708/hContainers/service"
	"os"
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
