package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"k3y0708/hContainers/colors"
	"k3y0708/hContainers/global"
	. "k3y0708/hContainers/types"
	. "k3y0708/hContainers/util"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

const (
	HCLOUD_API_URL = "https://api.hetzner.cloud/v1"
	CONFIG_FILE    = "config.ini"
)

//go:embed resources/cloudinit.yml
var cloudinit string

func main() {
	args := os.Args[1:]

	global.Client = hcloud.NewClient(hcloud.WithToken(os.Getenv("HCLOUD_TOKEN")))
	sshKeyPath, present := os.LookupEnv("HCONTAINERS_SSH_KEY_PATH")
	if !present {
		fmt.Println("HCONTAINERS_SSH_KEY_PATH not set")
		os.Exit(1)
	}
	sshKeyPath = strings.Replace(sshKeyPath, "~", os.Getenv("HOME"), 1)
	publicKeyBytes, err := os.ReadFile(sshKeyPath + ".pub")
	CheckError(err, "Failed to read public key", 1)
	global.PublicKey = string(publicKeyBytes)
	privateKeyBytes, err := os.ReadFile(sshKeyPath)
	CheckError(err, "Failed to read private key", 1)
	global.PrivateKey = string(privateKeyBytes)

	if len(args) == 0 {
		fmt.Println("No arguments provided")
		ShowHelp()
		os.Exit(1)
	}

	switch args[0] {
	case "help":
		ShowHelp()
	case "runner":
		switch args[1] {
		case "list":
			fmt.Printf("Available runners (Amount: %d):\n", GetAmountOfRunners())
			for _, server := range GetAllServerNames() {
				fmt.Printf("- %s\n", server)
			}
		case "create":
			if len(args) < 3 {
				fmt.Println("No runner name provided")
				os.Exit(1)
			}
			if !CheckIfRunnerNameIsValid(args[2]) {
				fmt.Println("Invalid runner name")
				os.Exit(1)
			}
			if CheckIfServerExists(args[2]) {
				fmt.Println("Runner already exists")
				os.Exit(1)
			}

			fmt.Println("Creating runner...")
			create_options := hcloud.ServerCreateOpts{
				Name:       args[2],
				ServerType: &hcloud.ServerType{Name: "cx11"},
				Image:      &hcloud.Image{Name: "ubuntu-22.04"},
				Location:   &hcloud.Location{City: "fsn1"},
				UserData:   strings.Replace(cloudinit, "{{{PUBLIC_KEY}}}", global.PublicKey, 1),
				PublicNet:  &hcloud.ServerCreatePublicNet{EnableIPv4: true, EnableIPv6: true},
				Labels:     map[string]string{"runner": "true"},
			}
			_, _, err := global.Client.Server.Create(context.Background(), create_options)
			CheckError(err, "Failed to create runner", 1)
			fmt.Println("Runner created")
		case "delete":
			if len(args) < 3 {
				fmt.Println("No runner name provided")
				os.Exit(1)
			}
			if !CheckIfServerExists(args[2]) {
				fmt.Println("Runner does not exist")
				os.Exit(1)
			}
			fmt.Println("Deleting runner...")
			_, _, err := global.Client.Server.DeleteWithResult(context.Background(), GetServerByName(args[2]))
			CheckError(err, "Failed to delete runner", 1)
			fmt.Println("Runner deleted")
		case "restart":
			if len(args) < 3 {
				fmt.Println("No runner name provided")
				os.Exit(1)
			}
			if !CheckIfServerExists(args[2]) {
				fmt.Println("Runner does not exist")
				os.Exit(1)
			}
			fmt.Println("Restarting runner...")
			_, _, err := global.Client.Server.Reboot(context.Background(), GetServerByName(args[2]))
			CheckError(err, "Failed to restart runner", 1)
			fmt.Println("Runner restarted")
		}
	case "container":
		switch args[1] {
		case "list":
			containers := GetAllContainers()
			fmt.Printf("Available containers (Amount: %d):\n", len(containers))
			var runningContainers, pausedContainers, exitedContainers, createdContainers, unknownContainers []Container
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
		case "create":
			if len(args) < 3 {
				fmt.Println("No runner name provided")
				os.Exit(1)
			}
			if !CheckIfServerExists(args[2]) {
				fmt.Println("Runner does not exist")
				os.Exit(1)
			}
			if len(args) < 4 {
				fmt.Println("No container name provided")
				os.Exit(1)
			}
			if !CheckIfContainerNameIsValid(args[3]) {
				fmt.Println("Invalid container name")
				os.Exit(1)
			}
			if CheckIfContainerExists(args[3]) {
				fmt.Println("Container already exists")
				os.Exit(1)
			}
			if len(args) < 5 {
				fmt.Println("No image name provided")
				os.Exit(1)
			}
			fmt.Println("Creating container...")
			if len(args) < 6 {
				_, err := RemoteRun(GetIpByName(args[2]), fmt.Sprintf("sudo nerdctl run -d --restart=unless-stopped --name %s %s", args[3], args[4]))
				CheckError(err, "Failed to create container", 1)
				fmt.Println("Container created")
			} else {
				_, err := RemoteRun(GetIpByName(args[2]), fmt.Sprintf("sudo nerdctl run -d --restart=unless-stopped -p %s --name %s %s", args[5], args[3], args[4]))
				CheckError(err, "Failed to create container", 1)
				fmt.Printf("Container created (with port mapping %s)\n", args[5])
			}
		case "delete":
			if len(args) < 3 {
				fmt.Println("No container name provided")
				os.Exit(1)
			}
			if !CheckIfContainerExists(args[2]) {
				fmt.Println("Container does not exist")
				os.Exit(1)
			}
			fmt.Println("Deleting container...")
			container := GetContainerByName(args[2])
			_, err := RemoteRun(GetIpByName(container.Runner), fmt.Sprintf("sudo nerdctl rm %s -f", container.Name))
			CheckError(err, "Failed to delete container", 1)
			fmt.Printf("Container %s deleted", container.Name)
		case "exec":
			if len(args) < 3 {
				fmt.Println("No container name provided")
				os.Exit(1)
			}
			if !CheckIfContainerExists(args[2]) {
				fmt.Println("Container does not exist")
				os.Exit(1)
			}
			if len(args) < 4 {
				fmt.Println("No command provided")
				os.Exit(1)
			}
			container := GetContainerByName(args[2])
			fmt.Println("Note: Commands which require user input are not supported")
			fmt.Printf("Connecting to %s...\n", container.Name)
			stdout, err := RemoteRun(GetIpByName(container.Runner), fmt.Sprintf("sudo nerdctl exec %s %s", container.Name, strings.Join(args[3:], " ")))
			CheckError(err, "Failed to connect to container", 1)
			fmt.Print(stdout)
			fmt.Println("Disconnected from container")
		case "stop":
			if len(args) < 3 {
				fmt.Println("No container name provided")
				os.Exit(1)
			}
			if !CheckIfContainerExists(args[2]) {
				fmt.Println("Container does not exist")
				os.Exit(1)
			}
			container := GetContainerByName(args[2])
			fmt.Println("Stopping container...")
			_, err := RemoteRun(GetIpByName(container.Runner), fmt.Sprintf("sudo nerdctl stop %s", container.Name))
			CheckError(err, "Failed to stop container", 1)
			fmt.Printf("Container %s stopped\n", container.Name)
		case "start":
			if len(args) < 3 {
				fmt.Println("No container name provided")
				os.Exit(1)
			}
			if !CheckIfContainerExists(args[2]) {
				fmt.Println("Container does not exist")
				os.Exit(1)
			}
			container := GetContainerByName(args[2])
			fmt.Println("Starting container...")
			_, err := RemoteRun(GetIpByName(container.Runner), fmt.Sprintf("sudo nerdctl start %s", container.Name))
			CheckError(err, "Failed to start container", 1)
			fmt.Printf("Container %s started\n", container.Name)
		case "restart":
			if len(args) < 3 {
				fmt.Println("No container name provided")
				os.Exit(1)
			}
			if !CheckIfContainerExists(args[2]) {
				fmt.Println("Container does not exist")
				os.Exit(1)
			}
			container := GetContainerByName(args[2])
			fmt.Println("Restarting container...")
			_, err := RemoteRun(GetIpByName(container.Runner), fmt.Sprintf("sudo nerdctl restart %s", container.Name))
			CheckError(err, "Failed to restart container", 1)
			fmt.Printf("Container %s restarted\n", container.Name)
		case "pause":
			if len(args) < 3 {
				fmt.Println("No container name provided")
				os.Exit(1)
			}
			if !CheckIfContainerExists(args[2]) {
				fmt.Println("Container does not exist")
				os.Exit(1)
			}
			container := GetContainerByName(args[2])
			fmt.Println("Pausing container...")
			_, err := RemoteRun(GetIpByName(container.Runner), fmt.Sprintf("sudo nerdctl pause %s", container.Name))
			CheckError(err, "Failed to pause container", 1)
			fmt.Printf("Container %s paused\n", container.Name)
		case "unpause":
			if len(args) < 3 {
				fmt.Println("No container name provided")
				os.Exit(1)
			}
			if !CheckIfContainerExists(args[2]) {
				fmt.Println("Container does not exist")
				os.Exit(1)
			}
			container := GetContainerByName(args[2])
			fmt.Println("Unpausing container...")
			_, err := RemoteRun(GetIpByName(container.Runner), fmt.Sprintf("sudo nerdctl unpause %s", container.Name))
			CheckError(err, "Failed to unpause container", 1)
			fmt.Printf("Container %s unpaused\n", container.Name)
		case "logs":
			if len(args) < 3 {
				fmt.Println("No container name provided")
				os.Exit(1)
			}
			if !CheckIfContainerExists(args[2]) {
				fmt.Println("Container does not exist")
				os.Exit(1)
			}
			container := GetContainerByName(args[2])
			fmt.Println("Getting logs...")
			stdout, err := RemoteRun(GetIpByName(container.Runner), fmt.Sprintf("sudo nerdctl logs %s", container.Name))
			CheckError(err, "Failed to get logs", 1)
			fmt.Print(stdout)
			fmt.Println("Disconnected from container")
		}
	}
}
