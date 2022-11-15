package assembler

import (
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
