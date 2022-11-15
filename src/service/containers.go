package service

import (
	"fmt"
	"k3y0708/hContainers/assembler"
	"k3y0708/hContainers/colors"
	"k3y0708/hContainers/global"
	. "k3y0708/hContainers/types"
	"k3y0708/hContainers/util"
	"os"
	"regexp"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func ContainerList() {
	var containers = getAllContainers()
	fmt.Printf("Available containers (Amount: %d):\n", len(containers))
	runningContainers, pausedContainers, exitedContainers, createdContainers, unknownContainers := assembler.ContainerList(containers)
	if len(runningContainers) > 0 {
		fmt.Printf("[%sRunning%s]\n", colors.GREEN, colors.RESET)
		for _, container := range runningContainers {
			fmt.Printf("  - %s (%s:%s)\n", container.Name, container.Image, container.Version)
		}
	}
	if len(pausedContainers) > 0 {
		fmt.Printf("[%sPaused%s]\n", colors.YELLOW, colors.RESET)
		for _, container := range pausedContainers {
			fmt.Printf("  - %s (%s:%s)\n", container.Name, container.Image, container.Version)
		}
	}
	if len(exitedContainers) > 0 {
		fmt.Printf("[%sStopped%s]\n", colors.RED, colors.RESET)
		for _, container := range exitedContainers {
			fmt.Printf("  - %s (%s:%s)\n", container.Name, container.Image, container.Version)
		}
	}
	if len(createdContainers) > 0 {
		fmt.Printf("[%sCreated%s]\n", colors.BLUE, colors.RESET)
		for _, container := range createdContainers {
			fmt.Printf("  - %s (%s:%s)\n", container.Name, container.Image, container.Version)
		}
	}
	if len(unknownContainers) > 0 {
		fmt.Printf("[%sOther%s]\n", colors.GREY, colors.RESET)
		for _, container := range unknownContainers {
			fmt.Printf("  - %s (Image: %s:%s Status: %s)\n", container.Name, container.Image, container.Version, container.Status)
		}
	}
}

func ContainerCreate(runnerName string, name string, image string, flags string) {
	if !checkIfServerExists(runnerName) {
		fmt.Println("Runner does not exist")
		os.Exit(1)
	}
	if checkIfContainerExists(name) {
		fmt.Println("Container already exists")
		os.Exit(1)
	}
	if !checkIfContainerNameIsValid(name) {
		fmt.Println("Container name is not valid")
		os.Exit(1)
	}
	fmt.Println("Creating container...")
	_, err := util.RemoteRun(util.GetIpByName(runnerName), fmt.Sprintf(global.CREATE, flags, name, image))
	util.CheckError(err, "Failed to create container", 1)
	fmt.Println("Container created")
}

func ContainerDelete(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	fmt.Println("Deleting container...")
	container := util.GetContainerByName(name)
	_, err := util.RemoteRun(util.GetIpByName(container.Runner), fmt.Sprintf(global.DELETE, container.Name))
	util.CheckError(err, "Failed to delete container", 1)
	fmt.Printf("Container %s deleted\n", container.Name)
}

func ContainerStart(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	fmt.Println("Starting container...")
	container := util.GetContainerByName(name)
	_, err := util.RemoteRun(util.GetIpByName(container.Runner), fmt.Sprintf(global.START, container.Name))
	util.CheckError(err, "Failed to start container", 1)
	fmt.Printf("Container %s started\n", container.Name)
}

func ContainerStop(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	fmt.Println("Stopping container...")
	container := util.GetContainerByName(name)
	_, err := util.RemoteRun(util.GetIpByName(container.Runner), fmt.Sprintf(global.STOP, container.Name))
	util.CheckError(err, "Failed to stop container", 1)
	fmt.Printf("Container %s stopped\n", container.Name)
}

func ContainerRestart(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	fmt.Println("Restarting container...")
	container := util.GetContainerByName(name)
	_, err := util.RemoteRun(util.GetIpByName(container.Runner), fmt.Sprintf(global.RESTART, container.Name))
	util.CheckError(err, "Failed to restart container", 1)
	fmt.Printf("Container %s restarted\n", container.Name)
}

func ContainerPause(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	fmt.Println("Pausing container...")
	container := util.GetContainerByName(name)
	_, err := util.RemoteRun(util.GetIpByName(container.Runner), fmt.Sprintf(global.PAUSE, container.Name))
	util.CheckError(err, "Failed to pause container", 1)
	fmt.Printf("Container %s paused\n", container.Name)
}

func ContainerUnpause(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	fmt.Println("Unpausing container...")
	container := util.GetContainerByName(name)
	_, err := util.RemoteRun(util.GetIpByName(container.Runner), fmt.Sprintf(global.UNPAUSE, container.Name))
	util.CheckError(err, "Failed to unpause container", 1)
	fmt.Printf("Container %s unpaused\n", container.Name)
}

func ContainerExec(name string, command string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	fmt.Println(colors.GREY + "Note: Commands which require user input are not supported" + colors.RESET)
	fmt.Println("Executing command...")
	container := util.GetContainerByName(name)
	stdout, err := util.RemoteRun(util.GetIpByName(container.Runner), fmt.Sprintf(global.EXEC, container.Name, command))
	util.CheckError(err, "Failed to exec command", 1)
	fmt.Println(stdout)
	fmt.Println("Executed command")
}

func ContainerLogs(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	fmt.Println("Getting logs...")
	container := util.GetContainerByName(name)
	stdout, err := util.RemoteRun(util.GetIpByName(container.Runner), fmt.Sprintf(global.LOGS, container.Name))
	util.CheckError(err, "Failed to get logs", 1)
	fmt.Println(stdout)
	fmt.Println("Got logs")
}

func checkIfContainerExists(containerName string) bool {
	containers := getAllContainers()
	for _, container := range containers {
		if container.Name == containerName+"-001" {
			return true
		}
	}
	return false
}

func checkIfContainerNameIsValid(containerName string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9]+(?:[._-](?:[A-Za-z0-9]+))*$`).MatchString(containerName)
}

func getAllContainers() []Container {
	var containers []Container
	for _, server := range GetAllServers() {
		containers = append(containers, getContainersFromRunner(server)...)
	}
	return containers
}

func getContainersFromRunner(runner *hcloud.Server) []Container {
	var containers []Container

	stdout, err := util.RemoteRun(runner.PublicNet.IPv4.IP.String(), global.LIST)
	util.CheckError(err, "Failed to get containers", -1)

	for _, container := range strings.Split(stdout, "\n") {
		if strings.TrimSpace(container) != "" {
			elements := strings.Split(container, " ")
			containers = append(containers,
				Container{
					Name:     strings.Split(elements[0], "-")[0],
					Instance: strings.Split(elements[0], "-")[1],
					ID:       elements[1],
					Image:    strings.Replace(strings.Split(elements[2], ":")[0], "docker.io/library/", "", 1),
					Version:  strings.Split(elements[2], ":")[1],
					Runner:   runner.Name,
					Status:   strings.Split(elements[3], " ")[0],
				})
		}
	}
	return containers
}
