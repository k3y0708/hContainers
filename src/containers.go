package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hContainers/hContainers/assembler"
	"github.com/hContainers/hContainers/printer"
	"github.com/hContainers/hContainers/service"
	"github.com/hContainers/hContainers/util"
)

func cliContainers(args []string) {
	if len(args) == 0 {
		cliContainersHelp()
		os.Exit(1)
	}
	switch args[0] {
	case "list":
		printer.ContainerInstancesMap(assembler.ContainersToContainerInstances(service.GetAllContainers()))
	case "create":
		util.CheckLength(args, 2, "No runner name provided", 1)
		util.CheckLength(args, 3, "No container name provided", 1)
		util.CheckLength(args, 4, "No container image provided", 1)
		fmt.Println("Creating container...")
		service.ContainerCreate(args[1], args[2], service.FindLowestPortPrefix(), "001", args[3])
		fmt.Println("Container created")
	case "delete":
		util.CheckLength(args, 2, "No container name provided", 1)
		fmt.Println("Deleting container...")
		service.ContainerDeleteAll(args[1])
		fmt.Printf("Container %s deleted\n", args[1])
	case "exec":
		util.CheckLength(args, 2, "No container name provided", 1)
		util.CheckLength(args, 3, "No instance provided", 1)
		util.CheckLength(args, 3, "No command provided", 1)
		service.ContainerExec(args[1], args[2], strings.Join(args[3:], " "))
		fmt.Println("Executed command")
	case "stop":
		util.CheckLength(args, 2, "No container name provided", 1)
		service.ContainerStop(args[1])
		fmt.Printf("Container %s stopped\n", args[1])
	case "start":
		util.CheckLength(args, 2, "No container name provided", 1)
		service.ContainerStart(args[1])
		fmt.Printf("Container %s started\n", args[1])
	case "restart":
		util.CheckLength(args, 2, "No container name provided", 1)
		service.ContainerRestart(args[1])
		fmt.Printf("Container %s restarted\n", args[1])
	case "pause":
		util.CheckLength(args, 2, "No container name provided", 1)
		service.ContainerPause(args[1])
		fmt.Printf("Container %s paused\n", args[1])
	case "unpause":
		util.CheckLength(args, 2, "No container name provided", 1)
		service.ContainerUnpause(args[1])
		fmt.Printf("Container %s unpaused\n", args[1])
	case "logs":
		util.CheckLength(args, 2, "No container name provided", 1)
		util.CheckLength(args, 3, "No instance provided", 1)
		service.ContainerLogs(args[1], args[2])
	case "scale":
		util.CheckLength(args, 2, "No container name provided", 1)
		util.CheckLength(args, 3, "No scale provided", 1)
		scale, err := strconv.Atoi(args[2])
		util.CheckError(err, "Invalid scale", 1)
		fmt.Println("Scaling container...")
		service.ContainerScale(args[1], scale)
		fmt.Printf("Container %s scaled to %d instances\n", args[1], scale)
	default:
		fmt.Println("Invalid command")
		nearestCommand := util.FindNearestCommand(args[0], []string{"list", "create", "delete", "exec", "stop", "start", "restart", "pause", "unpause", "logs", "scale"})
		if nearestCommand != "" {
			fmt.Println("Did you mean 'hContainers containers " + nearestCommand + "'?")
		} else {
			cliContainersHelp()
		}
	}
}
