package service

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/hContainers/hContainers/colors"
	"github.com/hContainers/hContainers/global"
	. "github.com/hContainers/hContainers/types"
	"github.com/hContainers/hContainers/util"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func ContainerCreate(runnerName, name, portPrefix, instance, image string) {
	if !checkIfServerExists(runnerName) {
		fmt.Printf("Runner %s does not exist\n", runnerName)
		os.Exit(1)
	}
	if checkIfContainerInstanceExists(name, instance) {
		fmt.Println("Container already exists")
		os.Exit(1)
	}
	if !checkIfContainerNameIsValid(name) {
		fmt.Println("Container name is not valid")
		os.Exit(1)
	}
	_, err := util.RemoteRun(util.GetIpByName(runnerName), fmt.Sprintf(global.CREATE, "-p 8080:"+portPrefix+instance, name, portPrefix, instance, image), 2)
	util.CheckError(err, "Failed to create container", 1)
}

func ContainerDeleteAll(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	container := getContainersByName(name)
	for _, c := range container {
		containerDelete(c.Name, c.Instance)
	}
}

func containerDelete(name, instance string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	container := getContainerByNameAndInstance(name, instance)
	_, err := util.RemoteRun(util.GetIpByName(container.Runner), fmt.Sprintf(global.DELETE, container.Name, container.PortPrefix, container.Instance), 2)
	if err != nil {
		fmt.Println("Failed to delete container")
		fmt.Println(err)
	}
	util.CheckError(err, "Failed to delete container", 1)
}

func ContainerStart(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	containers := getContainersByName(name)
	for _, c := range containers {
		_, err := util.RemoteRun(util.GetIpByName(c.Runner), fmt.Sprintf(global.START, c.Name, c.PortPrefix, c.Instance), 2)
		util.CheckError(err, "Failed to start container", 1)
	}
}

func ContainerStop(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	containers := getContainersByName(name)
	for _, c := range containers {
		_, err := util.RemoteRun(util.GetIpByName(c.Runner), fmt.Sprintf(global.STOP, c.Name, c.PortPrefix, c.Instance), 2)
		util.CheckError(err, "Failed to stop container", 1)
	}
}

func ContainerRestart(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	containers := getContainersByName(name)
	for _, c := range containers {
		_, err := util.RemoteRun(util.GetIpByName(c.Runner), fmt.Sprintf(global.RESTART, c.Name, c.PortPrefix, c.Instance), 2)
		util.CheckError(err, "Failed to restart container", 1)
	}
}

func ContainerPause(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	fmt.Println("Pausing container...")
	containers := getContainersByName(name)
	for _, c := range containers {
		_, err := util.RemoteRun(util.GetIpByName(c.Runner), fmt.Sprintf(global.PAUSE, c.Name, c.PortPrefix, c.Instance), 2)
		util.CheckError(err, "Failed to pause container", 1)
	}
}

func ContainerUnpause(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	fmt.Println("Unpausing container...")
	containers := getContainersByName(name)
	for _, c := range containers {
		_, err := util.RemoteRun(util.GetIpByName(c.Runner), fmt.Sprintf(global.UNPAUSE, c.Name, c.PortPrefix, c.Instance), 2)
		util.CheckError(err, "Failed to unpause container", 1)
	}
}

func ContainerExec(name, instance, command string) {
	if !checkIfContainerInstanceExists(name, instance) {
		fmt.Println("Container Instance does not exist")
		os.Exit(1)
	}
	fmt.Println(colors.GREY + "Note: Commands which require user input are not supported" + colors.RESET)
	container := getContainerByNameAndInstance(name, instance)
	stdout, err := util.RemoteRun(util.GetIpByName(container.Runner), fmt.Sprintf(global.EXEC, container.Name, container.PortPrefix, container.Instance, command), 2)
	util.CheckError(err, "Failed to exec command", 1)
	fmt.Println(stdout)
}

func ContainerLogs(name, instance string) {
	if !checkIfContainerInstanceExists(name, instance) {
		fmt.Println("Container Instance does not exist")
		os.Exit(1)
	}
	container := getContainerByNameAndInstance(name, instance)
	stdout, err := util.RemoteRun(util.GetIpByName(container.Runner), fmt.Sprintf(global.LOGS, container.Name, container.PortPrefix, container.Instance), 2)
	util.CheckError(err, "Failed to get logs", 1)
	fmt.Println(stdout)
}

func ContainerScale(name string, amount int) {
	if amount < 1 {
		fmt.Println("Scale must be greater than 0")
		os.Exit(1)
	}
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	container := getContainersByName(name)
	if len(container) == amount {
		fmt.Println("Container already scaled to that amount")
		os.Exit(0)
	}
	if len(container) > amount {
		containerScaleDown(name, len(container)-amount)
	} else {
		containerScaleUp(name, amount-len(container))
	}
}

func checkIfContainerExists(containerName string) bool {
	containers := GetAllContainers()
	for _, container := range containers {
		if container.Name == containerName {
			return true
		}
	}
	return false
}

func checkIfContainerInstanceExists(containerName string, instance string) bool {
	containers := GetAllContainers()
	for _, container := range containers {
		if container.Name == containerName && container.Instance == instance {
			return true
		}
	}
	return false
}

func checkIfContainerNameIsValid(containerName string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9]+(?:[._](?:[A-Za-z0-9]+))*$`).MatchString(containerName)
}

func GetAllContainers() []Container {
	var containers []Container
	for _, server := range GetAllServers() {
		containers = append(containers, getContainersFromRunner(server)...)
	}
	return containers
}

func getContainersFromRunner(runner *hcloud.Server) []Container {
	var containers []Container

	stdout, err := util.RemoteRun(runner.PublicNet.IPv4.IP.String(), global.LIST, 2)
	util.CheckError(err, "Failed to get containers", -1)

	for _, container := range strings.Split(stdout, "\n") {
		if strings.TrimSpace(container) != "" {
			elements := strings.Split(container, " ")
			containers = append(containers,
				Container{
					Name:       strings.Split(elements[0], "-")[0],
					PortPrefix: strings.Split(elements[0], "-")[1],
					Instance:   strings.Split(elements[0], "-")[2],
					ID:         elements[1],
					Image:      strings.Replace(strings.Split(elements[2], ":")[0], "docker.io/library/", "", 1),
					Version:    strings.Split(elements[2], ":")[1],
					Runner:     runner.Name,
					Status:     strings.Split(elements[3], " ")[0],
				})
		}
	}
	return containers
}

func getContainerByName(containerName string) Container {
	containers := GetAllContainers()
	for _, container := range containers {
		if container.Name == containerName {
			return container
		}
	}
	return Container{}
}

func getContainerByNameAndInstance(containerName string, instance string) Container {
	containers := GetAllContainers()
	for _, container := range containers {
		if container.Name == containerName && container.Instance == instance {
			return container
		}
	}
	return Container{}
}

/*
Returns all instances of a container

@param containerName: Name of the container
@return: List of containers
*/
func getContainersByName(containerName string) []Container {
	var containers []Container
	for _, container := range GetAllContainers() {
		if container.Name == containerName {
			containers = append(containers, container)
		}
	}
	return containers
}

func getFreeInstanceNumber(containerName string) string {
	containers := getContainersByName(containerName)
	if len(containers) == 0 {
		return "001"
	}
	var instanceNumbers []int
	for _, container := range containers {
		instanceNumber, _ := strconv.Atoi(container.Instance)
		instanceNumbers = append(instanceNumbers, instanceNumber)
	}
	sort.Ints(instanceNumbers)
	for i := 0; i < len(instanceNumbers); i++ {
		if instanceNumbers[i] != i+1 {
			return fmt.Sprintf("%03d", i+1)
		}
	}
	return fmt.Sprintf("%03d", len(instanceNumbers)+1)
}

func containerScaleUp(containerName string, amount int) {
	container := getContainerByName(containerName)
	for i := 0; i < amount; i++ {
		instanceNumber := getFreeInstanceNumber(containerName)
		ContainerCreate(container.Runner, container.Name, container.PortPrefix, instanceNumber, container.GetFullImageName())
	}
}

func containerScaleDown(containerName string, amount int) {
	containers := getContainersByName(containerName)
	sort.Slice(containers, func(i, j int) bool {
		return containers[i].Instance < containers[j].Instance
	})
	for i := 0; i < amount; i++ {
		containerDelete(containers[i].Name, containers[i].Instance)
	}
}

func FindLowestPortPrefix() string {
	var portPrefixes []int
	var scannedContainers []string
	for _, container := range GetAllContainers() {
		if util.Contains(scannedContainers, container.Name) {
			continue
		}
		portPrefix, _ := strconv.Atoi(container.PortPrefix)
		portPrefixes = append(portPrefixes, portPrefix-1)
		scannedContainers = append(scannedContainers, container.Name)
	}
	sort.Ints(portPrefixes)
	for i := 0; i < len(portPrefixes); i++ {
		if portPrefixes[i] != i+1 {
			return fmt.Sprintf("%02d", i+1+1)
		}
	}
	return fmt.Sprintf("%02d", len(portPrefixes)+1+1)
}
