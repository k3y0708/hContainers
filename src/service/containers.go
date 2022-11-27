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
	"sort"
	"strconv"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func ContainerList() {
	var containers = getAllContainers()
	fmt.Printf("Available containers (Amount: %d):\n", len(containers))
	assembler.ContainerToTable(containers)
}

func ContainerCreate(runnerName string, name string, instance string, image string) {
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
	fmt.Println("Creating container...")
	_, err := util.RemoteRun(util.GetIpByName(runnerName), fmt.Sprintf(global.CREATE, "-p 8080:8"+instance, name, instance, image))
	util.CheckError(err, "Failed to create container", 1)
	fmt.Println("Container created")
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
	fmt.Printf("Container %s deleted\n", name)
}

func containerDelete(name, instance string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	container := getContainerByNameAndInstance(name, instance)
	_, err := util.RemoteRun(util.GetIpByName(container.Runner), fmt.Sprintf(global.DELETE, container.Name, container.Instance))
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
		_, err := util.RemoteRun(util.GetIpByName(c.Runner), fmt.Sprintf(global.START, c.Name, c.Instance))
		util.CheckError(err, "Failed to start container", 1)
	}
	fmt.Printf("Container %s started\n", name)
}

func ContainerStop(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	containers := getContainersByName(name)
	for _, c := range containers {
		_, err := util.RemoteRun(util.GetIpByName(c.Runner), fmt.Sprintf(global.STOP, c.Name, c.Instance))
		util.CheckError(err, "Failed to stop container", 1)
	}
	fmt.Printf("Container %s stopped\n", name)
}

func ContainerRestart(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	containers := getContainersByName(name)
	for _, c := range containers {
		_, err := util.RemoteRun(util.GetIpByName(c.Runner), fmt.Sprintf(global.RESTART, c.Name, c.Instance))
		util.CheckError(err, "Failed to restart container", 1)
	}
	fmt.Printf("Container %s restarted\n", name)
}

func ContainerPause(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	fmt.Println("Pausing container...")
	containers := getContainersByName(name)
	for _, c := range containers {
		_, err := util.RemoteRun(util.GetIpByName(c.Runner), fmt.Sprintf(global.PAUSE, c.Name, c.Instance))
		util.CheckError(err, "Failed to pause container", 1)
	}
	fmt.Printf("Container %s paused\n", name)
}

func ContainerUnpause(name string) {
	if !checkIfContainerExists(name) {
		fmt.Println("Container does not exist")
		os.Exit(1)
	}
	fmt.Println("Unpausing container...")
	containers := getContainersByName(name)
	for _, c := range containers {
		_, err := util.RemoteRun(util.GetIpByName(c.Runner), fmt.Sprintf(global.UNPAUSE, c.Name, c.Instance))
		util.CheckError(err, "Failed to unpause container", 1)
	}
	fmt.Printf("Container %s unpaused\n", name)
}

func ContainerExec(name, instance, command string) {
	if !checkIfContainerInstanceExists(name, instance) {
		fmt.Println("Container Instance does not exist")
		os.Exit(1)
	}
	fmt.Println(colors.GREY + "Note: Commands which require user input are not supported" + colors.RESET)
	container := getContainerByNameAndInstance(name, instance)
	stdout, err := util.RemoteRun(util.GetIpByName(container.Runner), fmt.Sprintf(global.EXEC, container.Name, container.Instance, command))
	util.CheckError(err, "Failed to exec command", 1)
	fmt.Println(stdout)
	fmt.Println("Executed command")
}

func ContainerLogs(name, instance string) {
	if !checkIfContainerInstanceExists(name, instance) {
		fmt.Println("Container Instance does not exist")
		os.Exit(1)
	}
	container := getContainerByNameAndInstance(name, instance)
	stdout, err := util.RemoteRun(util.GetIpByName(container.Runner), fmt.Sprintf(global.LOGS, container.Name))
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
	fmt.Println("Scaling container...")
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
	containers := getAllContainers()
	for _, container := range containers {
		if container.Name == containerName {
			return true
		}
	}
	return false
}

func checkIfContainerInstanceExists(containerName string, instance string) bool {
	containers := getAllContainers()
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

func getContainerByName(containerName string) Container {
	containers := getAllContainers()
	for _, container := range containers {
		if container.Name == containerName {
			return container
		}
	}
	return Container{}
}

func getContainerByNameAndInstance(containerName string, instance string) Container {
	containers := getAllContainers()
	for _, container := range containers {
		if container.Name == containerName && container.Instance == instance {
			return container
		}
	}
	return Container{}
}

func getContainersByName(containerName string) []Container {
	var containers []Container
	for _, container := range getAllContainers() {
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
	fmt.Printf("Scaling up container %s by %d\n", containerName, amount)
	container := getContainerByName(containerName)
	for i := 0; i < amount; i++ {
		instanceNumber := getFreeInstanceNumber(containerName)
		ContainerCreate(container.Runner, containerName, instanceNumber, container.Image+":"+container.Version)
	}
	fmt.Println("Scaled up container")
}

func containerScaleDown(containerName string, amount int) {
	fmt.Printf("Scaling down container %s by %d\n", containerName, amount)
	containers := getContainersByName(containerName)
	sort.Slice(containers, func(i, j int) bool {
		return containers[i].Instance < containers[j].Instance
	})
	for i := 0; i < amount; i++ {
		containerDelete(containers[i].Name, containers[i].Instance)
	}
	fmt.Println("Scaled down container")
}
