package main

import (
	"k3y0708/hContainers/service"
	"k3y0708/hContainers/util"
	"os"
	"strconv"
	"strings"
)

func cliContainers(args []string) {
	if len(args) == 0 {
		cliContainersHelp()
		os.Exit(1)
	}
	switch args[0] {
	case "list":
		service.ContainerList()
	case "create":
		util.CheckLength(args, 2, "No runner name provided", 1)
		util.CheckLength(args, 3, "No container name provided", 1)
		util.CheckLength(args, 4, "No container image provided", 1)
		service.ContainerCreate(args[1], args[2], service.FindLowestPortPrefix(), "001", args[3])
	case "delete":
		util.CheckLength(args, 2, "No container name provided", 1)
		service.ContainerDeleteAll(args[1])
	case "exec":
		util.CheckLength(args, 2, "No container name provided", 1)
		util.CheckLength(args, 3, "No instance provided", 1)
		util.CheckLength(args, 3, "No command provided", 1)
		service.ContainerExec(args[1], args[2], strings.Join(args[3:], " "))
	case "stop":
		util.CheckLength(args, 2, "No container name provided", 1)
		service.ContainerStop(args[1])
	case "start":
		util.CheckLength(args, 2, "No container name provided", 1)
		service.ContainerStart(args[1])
	case "restart":
		util.CheckLength(args, 2, "No container name provided", 1)
		service.ContainerRestart(args[1])
	case "pause":
		util.CheckLength(args, 2, "No container name provided", 1)
		service.ContainerPause(args[1])
	case "unpause":
		util.CheckLength(args, 2, "No container name provided", 1)
		service.ContainerUnpause(args[1])
	case "logs":
		util.CheckLength(args, 2, "No container name provided", 1)
		util.CheckLength(args, 3, "No instance provided", 1)
		service.ContainerLogs(args[1], args[2])
	case "scale":
		util.CheckLength(args, 2, "No container name provided", 1)
		util.CheckLength(args, 3, "No scale provided", 1)
		scale, err := strconv.Atoi(args[2])
		util.CheckError(err, "Invalid scale", 1)
		service.ContainerScale(args[1], scale)
	default:
		cliContainersHelp()
	}
}
