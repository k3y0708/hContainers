package main

import (
	"fmt"
	"k3y0708/hContainers/service"
	"k3y0708/hContainers/util"
	"os"
)

func cliRunner(args []string) {
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
		service.RunnerCreate(args[1])
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
