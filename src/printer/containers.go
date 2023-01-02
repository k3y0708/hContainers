package printer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hContainers/hContainers/types"
)

func ContainerInstancesMap(containers map[string]types.ContainerInstances) {
	longestContainerName := len("Name")
	longestRunnerName := len("Runner")
	longestImageName := len("Image")
	longestImageVersion := len("Version")
	longestCount := len("Count")
	longestHealthy := len("Healthy")

	for name, container := range containers {
		if len(name) > longestContainerName {
			longestContainerName = len(name)
		}
		if len(container.Runner) > longestRunnerName {
			longestRunnerName = len(container.Runner)
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

	fmt.Printf("%s %s %s %s Healthy/Count\n", "Runner"+strings.Repeat(" ", longestRunnerName-len("Runner")+2), "Name"+strings.Repeat(" ", longestContainerName-len("Name")+2), "Image"+strings.Repeat(" ", longestImageName-len("Image")+2), "Version"+strings.Repeat(" ", longestImageVersion-len("Version")+2))
	for name, container := range containers {
		containerName := name + strings.Repeat(" ", longestContainerName-len(name)+2)
		runnerName := container.Runner + strings.Repeat(" ", longestRunnerName-len(container.Runner)+2)
		imageName := container.Image + strings.Repeat(" ", longestImageName-len(container.Image)+2)
		imageVersion := container.Version + strings.Repeat(" ", longestImageVersion-len(container.Version)+2)
		count := strconv.Itoa(container.Count)
		healthy := strconv.Itoa(container.Healthy)

		fmt.Printf("%s %s %s %s %s/%s\n", runnerName, containerName, imageName, imageVersion, healthy, count)
	}
}
