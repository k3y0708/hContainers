package main

import (
	"fmt"
	"os"

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
		fmt.Printf("Available runners (Amount: %d):\n", util.GetAmountOfRunners())
		for _, server := range util.GetAllServerNames() {
			fmt.Printf("- %s\n", server)
		}
	case "create":
		util.CheckLength(args, 2, "No runner name provided", 1)
		if !util.CheckIfRunnerNameIsValid(args[1]) {
			fmt.Println("Invalid runner name")
			os.Exit(1)
		}
		if util.CheckIfServerExists(args[1]) {
			fmt.Println("Runner already exists")
			os.Exit(1)
		}
		var flag types.FlagsServer
		_, err := flags.ParseArgs(&flag, args[1:])
		util.CheckError(err, "Failed to parse flags", 1)
		service.RunnerCreate(args[1], flag)
	case "delete":
		util.CheckLength(args, 2, "No runner name provided", 1)
		if !util.CheckIfServerExists(args[1]) {
			fmt.Println("Runner does not exist")
			os.Exit(1)
		}
		service.RunnerDelete(args[1])
	case "restart":
		util.CheckLength(args, 2, "No runner name provided", 1)
		if !util.CheckIfServerExists(args[1]) {
			fmt.Println("Runner does not exist")
			os.Exit(1)
		}
		service.RunnerRestart(args[1])
	default:
		cliRunnerHelp()
	}
}
