package main

import (
	"fmt"
	"os"

	"github.com/hContainers/hContainers/assembler"
	"github.com/hContainers/hContainers/printer"
	"github.com/hContainers/hContainers/service"
	"github.com/hContainers/hContainers/types"
	"github.com/hContainers/hContainers/util"

	"github.com/jessevdk/go-flags"
)

func cliRunner(args []string) {
	if len(args) == 0 {
		cliRunnerHelp()
		os.Exit(1)
	}
	switch args[0] {
	case "list":
		printer.RunnerList(assembler.ServerListToNames(service.GetAllServers()))
	case "create":
		util.CheckLength(args, 2, "No runner name provided", 1)
		if !util.CheckIfRunnerNameIsValid(args[1]) {
			fmt.Println("Invalid runner name")
			os.Exit(1)
		}
		var flag types.FlagsServer
		_, err := flags.ParseArgs(&flag, args[1:])
		util.CheckError(err, "Failed to parse flags", 1)

		fmt.Println("Creating runner...")
		service.RunnerCreate(args[1], flag)
		fmt.Println("Runner created")
	case "delete":
		util.CheckLength(args, 2, "No runner name provided", 1)
		if !util.CheckIfServerExists(args[1]) {
			fmt.Println("Runner does not exist")
			os.Exit(1)
		}

		fmt.Println("Deleting runner...")
		service.RunnerDelete(args[1])
		fmt.Println("Runner deleted")
	case "restart":
		util.CheckLength(args, 2, "No runner name provided", 1)
		if !util.CheckIfServerExists(args[1]) {
			fmt.Println("Runner does not exist")
			os.Exit(1)
		}

		fmt.Println("Restarting runner...")
		service.RunnerRestart(args[1])
		fmt.Println("Runner restarted")
	default:
		fmt.Println("Invalid command")
		nearestCommand := util.FindNearestCommand(args[0], []string{"list", "create", "delete", "restart"})
		if nearestCommand != "" {
			fmt.Println("Did you mean 'hContainers runners " + nearestCommand + "'?")
		} else {
			cliRunnerHelp()
		}
	}
}
