package assembler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hContainers/hContainers/types"
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

func ContainerToMap(containers []types.Container) {
	containerMap := make(map[string]types.ContainerInstances)
	for _, container := range containers {
		if _, ok := containerMap[container.Name]; !ok {
			containerMap[container.Name] = types.ContainerInstances{
				Image:   container.Image,
				Version: container.Version,
				Count:   0,
				Healthy: 0,
			}
		}

		tmp := containerMap[container.Name]
		tmp.Count++
		if container.Status == "Up" {
			tmp.Healthy++
		}
		containerMap[container.Name] = tmp
	}

	longestContainerName := len("Name")
	longestImageName := len("Image")
	longestImageVersion := len("Version")
	longestCount := len("Count")
	longestHealthy := len("Healthy")
	for name, container := range containerMap {
		if len(name) > longestContainerName {
			longestContainerName = len(name)
		}
		if len(container.Image) > longestImageName {
			longestImageName = len(container.Image)
		}
		if len(container.Version) > longestImageVersion {
			longestImageVersion = len(container.Version)
		}
		if len(strconv.Itoa(container.Count)) > longestCount {
			longestCount = len(strconv.Itoa(container.Count))
		}
		if len(strconv.Itoa(container.Healthy)) > longestHealthy {
			longestHealthy = len(strconv.Itoa(container.Healthy))
		}
	}

	fmt.Printf("%s %s %s Healthy/Count\n", "Name"+strings.Repeat(" ", longestContainerName-len("Name")+2), "Image"+strings.Repeat(" ", longestImageName-len("Image")+2), "Version"+strings.Repeat(" ", longestImageVersion-len("Version")+2))
	for name, container := range containerMap {
		containerName := name + strings.Repeat(" ", longestContainerName-len(name)+2)
		imageName := container.Image + strings.Repeat(" ", longestImageName-len(container.Image)+2)
		imageVersion := container.Version + strings.Repeat(" ", longestImageVersion-len(container.Version)+2)
		count := strconv.Itoa(container.Count)
		healthy := strconv.Itoa(container.Healthy)

		fmt.Printf("%s %s %s %s/%s\n", containerName, imageName, imageVersion, healthy, count)
	}
}
