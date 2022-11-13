package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

const (
	HCLOUD_API_URL = "https://api.hetzner.cloud/v1"
	CONFIG_FILE    = "config.ini"
)

//go:embed resources/cloudinit.yml
var cloudinit string

var publicKey string

var client *hcloud.Client

func main() {
	args := os.Args[1:]

	client = hcloud.NewClient(hcloud.WithToken(os.Getenv("HCLOUD_TOKEN")))
	sshKeyPath, present := os.LookupEnv("HCONTAINERS_SSH_KEY_PATH")
	if !present {
		fmt.Println("HCONTAINERS_SSH_KEY_PATH not set")
		os.Exit(1)
	}
	sshKeyPath = strings.Replace(sshKeyPath, "~", os.Getenv("HOME"), 1)
	publicKeyBytes, err := os.ReadFile(sshKeyPath + ".pub")
	checkError(err, "Failed to read public key", 1)
	publicKey = string(publicKeyBytes)
	privateKeyBytes, err := os.ReadFile(sshKeyPath)
	checkError(err, "Failed to read private key", 1)
	privateKey = string(privateKeyBytes)

	if len(args) == 0 {
		fmt.Println("No arguments provided")
		showHelp()
		os.Exit(1)
	}

	switch args[0] {
	case "help":
		showHelp()
	case "runner":
		switch args[1] {
		case "list":
			fmt.Printf("Available runners (Amount: %d):\n", getAmountOfRunners())
			for _, server := range getAllServerNames() {
				fmt.Printf("- %s\n", server)
			}
		case "create":
			if len(args) < 3 {
				fmt.Println("No runner name provided")
				os.Exit(1)
			}
			if !checkIfRunnerNameIsValid(args[2]) {
				fmt.Println("Invalid runner name")
				os.Exit(1)
			}
			if checkIfServerExists(args[2]) {
				fmt.Println("Runner already exists")
				os.Exit(1)
			}

			fmt.Println("Creating runner...")
			create_options := hcloud.ServerCreateOpts{
				Name:       args[2],
				ServerType: &hcloud.ServerType{Name: "cx11"},
				Image:      &hcloud.Image{Name: "ubuntu-22.04"},
				Location:   &hcloud.Location{City: "fsn1"},
				UserData:   strings.Replace(cloudinit, "{{{PUBLIC_KEY}}}", publicKey, 1),
				PublicNet:  &hcloud.ServerCreatePublicNet{EnableIPv4: true, EnableIPv6: true},
				Labels:     map[string]string{"runner": "true"},
			}
			_, _, err := client.Server.Create(context.Background(), create_options)
			checkError(err, "Failed to create runner", 1)
			fmt.Println("Runner created")
		case "delete":
			if len(args) < 3 {
				fmt.Println("No runner name provided")
				os.Exit(1)
			}
			if !checkIfServerExists(args[2]) {
				fmt.Println("Runner does not exist")
				os.Exit(1)
			}
			fmt.Println("Deleting runner...")
			_, _, err := client.Server.DeleteWithResult(context.Background(), getServerByName(args[2]))
			checkError(err, "Failed to delete runner", 1)
			fmt.Println("Runner deleted")
		}
	case "container":
		switch args[1] {
		case "list":
			containers := getAllContainers()
			fmt.Printf("Available containers (Amount: %d):\n", len(containers))
			for _, container := range containers {
				fmt.Printf("- %s\n", container.Name)
			}
		case "create":
			if len(args) < 3 {
				fmt.Println("No runner name provided")
				os.Exit(1)
			}
			if !checkIfServerExists(args[2]) {
				fmt.Println("Runner does not exist")
				os.Exit(1)
			}
			if len(args) < 4 {
				fmt.Println("No container name provided")
				os.Exit(1)
			}
			if !checkIfContainerNameIsValid(args[3]) {
				fmt.Println("Invalid container name")
				os.Exit(1)
			}
			if checkIfContainerExists(args[3]) {
				fmt.Println("Container already exists")
				os.Exit(1)
			}
			if len(args) < 5 {
				fmt.Println("No image name provided")
				os.Exit(1)
			}
			fmt.Println("Creating container...")
			if len(args) < 6 {
				_, err := RemoteRun(getIpByName(args[2]), fmt.Sprintf("sudo nerdctl run -d --name %s %s", args[3], args[4]))
				checkError(err, "Failed to create container", 1)
				fmt.Println("Container created")
			} else {
				_, err := RemoteRun(getIpByName(args[2]), fmt.Sprintf("sudo nerdctl run -d -p %s --name %s %s", args[5], args[3], args[4]))
				checkError(err, "Failed to create container", 1)
				fmt.Printf("Container created (with port mapping %s)\n", args[5])
			}
		case "delete":
			if len(args) < 3 {
				fmt.Println("No container name provided")
				os.Exit(1)
			}
			if !checkIfContainerExists(args[2]) {
				fmt.Println("Container does not exist")
				os.Exit(1)
			}
			fmt.Println("Deleting container...")
			container := getContainerByName(args[2])
			_, err := RemoteRun(getIpByName(container.Runner), fmt.Sprintf("sudo nerdctl rm %s -f", container.Name))
			checkError(err, "Failed to delete container", 1)
			fmt.Printf("Container %s deleted", container.Name)
		case "exec":
			if len(args) < 3 {
				fmt.Println("No container name provided")
				os.Exit(1)
			}
			if !checkIfContainerExists(args[2]) {
				fmt.Println("Container does not exist")
				os.Exit(1)
			}
			if len(args) < 4 {
				fmt.Println("No command provided")
				os.Exit(1)
			}
			container := getContainerByName(args[2])
			fmt.Println("Note: Commands which require user input are not supported")
			fmt.Printf("Connecting to %s...\n", container.Name)
			stdout, err := RemoteRun(getIpByName(container.Runner), fmt.Sprintf("sudo nerdctl exec %s %s", container.Name, strings.Join(args[3:], " ")))
			checkError(err, "Failed to connect to container", 1)
			fmt.Print(stdout)
			fmt.Println("Disconnected from container")
		}
	}
}

func getAllServers() []*hcloud.Server {
	servers, _, err := client.Server.List(context.Background(), hcloud.ServerListOpts{})
	checkError(err, "Failed to get servers", -1)
	return servers
}

func getAllServerNames() []string {
	var serverNames []string
	for _, server := range getAllServers() {
		serverNames = append(serverNames, server.Name)
	}
	return serverNames
}

func getServerByName(serverName string) *hcloud.Server {
	servers, _, err := client.Server.List(context.Background(), hcloud.ServerListOpts{})
	checkError(err, "Failed to get servers", -1)
	for _, server := range servers {
		if server.Name == serverName {
			return server
		}
	}
	return nil
}

func getIpByName(serverName string) string {
	return getServerByName(serverName).PublicNet.IPv4.IP.String()
}

func checkIfServerExists(serverName string) bool {
	servers := getAllServerNames()

	for _, server := range servers {
		if server == serverName {
			return true
		}
	}
	return false
}

func checkError(err error, message string, exitCode int) {
	if err != nil {
		fmt.Println(message)
		fmt.Println(err)
		os.Exit(exitCode)
	}
}

func checkIfRunnerNameIsValid(runnerName string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9-.]*$`).MatchString(runnerName)
}

func checkIfContainerNameIsValid(containerName string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9]+(?:[._-](?:[A-Za-z0-9]+))*$`).MatchString(containerName)
}

func showHelp() {
	fmt.Println("Usage: hContainers <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  runner list                                                         - List all available runners")
	fmt.Println("  runner create <runner-name>                                         - Create a new runner")
	fmt.Println("  runner delete <runner-name>                                         - Delete a runner")
	fmt.Println("  container list                                                      - List all available containers")
	fmt.Println("  container create <runner-name> <container-name> <image>             - Create a new container")
	fmt.Println("  container create <runner-name> <container-name> <image> <port:port> - Create a new container with port mapping")
	fmt.Println("  container delete <name>                                             - Delete a container")
	fmt.Println("  container exec <name> <command>                                     - Execute a command in a container (Non-interactive)")
	fmt.Println("  help                                                                - Show this help message")
}

func getAmountOfRunners() int {
	servers, _, err := client.Server.List(context.Background(), hcloud.ServerListOpts{ListOpts: hcloud.ListOpts{LabelSelector: "runner"}})
	checkError(err, "Failed to get servers", -1)
	return len(servers)
}

func getAllContainers() []Container {
	var containers []Container
	for _, server := range getAllServers() {
		var stdout string
		var err error
		stdout, err = RemoteRun(server.PublicNet.IPv4.IP.String(), "sudo nerdctl container ls --format '{{.Names}} {{.ID}} {{.Image}}'")
		if err != nil {
			fmt.Println("An error occured while getting containers")
		}
		for _, container := range strings.Split(stdout, "\n") {
			if strings.TrimSpace(container) != "" {
				elements := strings.Split(container, " ")
				containers = append(containers, Container{Name: elements[0], ID: elements[1], Image: elements[2], Runner: server.Name})
			}
		}
	}
	return containers
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

func getContainerByName(containerName string) Container {
	containers := getAllContainers()
	for _, container := range containers {
		if container.Name == containerName {
			return container
		}
	}
	return Container{}
}
