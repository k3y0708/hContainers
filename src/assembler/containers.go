package assembler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hContainers/hContainers/types"
)

func ContainersToContainerInstances(containers []types.Container) map[string]types.ContainerInstances {
	containerMap := make(map[string]types.ContainerInstances)
	for _, container := range containers {
		if _, ok := containerMap[container.Name]; !ok {
			containerMap[container.Name] = types.ContainerInstances{
				Image:   container.Image,
				Runner:  container.Runner,
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
	return containerMap
}

func ContainerToMap(containers []types.Container) {
	containerMap := ContainersToContainerInstances(containers)

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
