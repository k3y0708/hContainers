package assembler

import (
	"fmt"
	"k3y0708/hContainers/types"
	"strings"
)

func ContainerList(containers []types.Container) ([]types.Container, []types.Container, []types.Container, []types.Container, []types.Container) {
	var runningContainers, pausedContainers, exitedContainers, createdContainers, unknownContainers []types.Container
	for _, container := range containers {
		switch strings.ToLower(container.Status) {
		case "up":
			runningContainers = append(runningContainers, container)
		case "paused":
			pausedContainers = append(pausedContainers, container)
		case "exited":
			exitedContainers = append(exitedContainers, container)
		case "created":
			createdContainers = append(createdContainers, container)
		default:
			unknownContainers = append(unknownContainers, container)
		}
	}
	return runningContainers, pausedContainers, exitedContainers, createdContainers, unknownContainers
}

func ContainerToTable(containers []types.Container) {
	longestContainerName := 0
	longestContainerInstance := 0
	longestImageName := 0
	longestImageVersion := 0
	longestStatus := 0
	for _, container := range containers {
		if len(container.Name) > longestContainerName {
			longestContainerName = len(container.Name)
		}
		if len(container.Instance) > longestContainerInstance {
			longestContainerInstance = len(container.Instance)
		}
		if len(container.Image) > longestImageName {
			longestImageName = len(container.Image)
		}
		if len(container.Version) > longestImageVersion {
			longestImageVersion = len(container.Version)
		}
		if len(container.Status) > longestStatus {
			longestStatus = len(container.Status)
		}
	}

	lastContainerName := ""
	for _, container := range containers {
		var containerName, imageName, imageVersion string
		if container.Name == lastContainerName {
			containerName = strings.Repeat(" ", longestContainerName+2)
			imageName = strings.Repeat(" ", longestImageName+2)
			imageVersion = strings.Repeat(" ", longestImageVersion+2)
		} else {
			containerName = container.Name + strings.Repeat(" ", longestContainerName-len(container.Name)+2)
			imageName = container.Image + strings.Repeat(" ", longestImageName-len(container.Image)+2)
			imageVersion = container.Version + strings.Repeat(" ", longestImageVersion-len(container.Version)+2)
			lastContainerName = container.Name
		}
		containerInstance := container.Instance + strings.Repeat(" ", longestContainerInstance-len(container.Instance)+2)
		status := container.Status + strings.Repeat(" ", longestStatus-len(container.Status)+2)

		fmt.Printf("%s %s %s %s %s\n", containerName, imageName, imageVersion, containerInstance, status)
	}
}
