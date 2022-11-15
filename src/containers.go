package main

import (
	"k3y0708/hContainers/service"
	"k3y0708/hContainers/util"
	"strings"
)

func cliContainers(args []string) {
	switch args[0] {
	case "list":
		service.ContainerList()
	case "create":
		util.CheckLength(args, 2, "No runner name provided", 1)
		util.CheckLength(args, 3, "No container name provided", 1)
		util.CheckLength(args, 4, "No container image provided", 1)
		service.ContainerCreate(args[1], args[2], args[3], strings.Join(args[4:], " "))
	case "delete":
		util.CheckLength(args, 2, "No container name provided", 1)
		service.ContainerDelete(args[1])
	case "exec":
		util.CheckLength(args, 2, "No container name provided", 1)
		util.CheckLength(args, 3, "No command provided", 1)
		service.ContainerExec(args[1], strings.Join(args[2:], " "))
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
		service.ContainerLogs(args[1])
	default:
		cliContainersHelp()
	}
}
