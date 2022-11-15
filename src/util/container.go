package util

import (
	"fmt"
	"strings"

	. "k3y0708/hContainers/types"
)

func GetAllContainersOld() []Container {
	var containers []Container
	for _, server := range GetAllServers() {
		var stdout string
		var err error
		stdout, err = RemoteRun(server.PublicNet.IPv4.IP.String(), "sudo nerdctl container ls --all --format '{{.Names}} {{.ID}} {{.Image}} {{.Status}}'")
		if err != nil {
			fmt.Println("An error occured while getting containers")
		}
		for _, container := range strings.Split(stdout, "\n") {
			if strings.TrimSpace(container) != "" {
				elements := strings.Split(container, " ")
				containers = append(containers,
					Container{
						Name:    elements[0],
						ID:      elements[1],
						Image:   strings.Replace(strings.Split(elements[2], ":")[0], "docker.io/library/", "", 1),
						Version: strings.Split(elements[2], ":")[1],
						Runner:  server.Name,
						Status:  strings.Split(elements[3], " ")[0],
					})
			}
		}
	}
	return containers
}

func CheckIfContainerExists(containerName string) bool {
	containers := GetAllContainersOld()
	for _, container := range containers {
		if container.Name == containerName+"-001" {
			return true
		}
	}
	return false
}

func GetContainerByName(containerName string) Container {
	containers := GetAllContainersOld()
	for _, container := range containers {
		if container.Name == containerName+"-001" {
			return container
		}
	}
	return Container{}
}
