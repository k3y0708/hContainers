package main

import (
	"fmt"
	"os"
	"strings"

	"k3y0708/hContainers/global"
	. "k3y0708/hContainers/util"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

const (
	HCLOUD_API_URL = "https://api.hetzner.cloud/v1"
	CONFIG_FILE    = "config.ini"
)

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
		cliHelp()
		os.Exit(1)
	}

	switch args[0] {
	case "help":
		cliHelp()
	case "runner":
		cliRunner(args[1:])
	case "containers":
		cliContainers(args[1:])
	case "container":
		fmt.Println("Invalid command\nDid you mean 'containers'?")
	default:
		fmt.Println("Invalid command")
		cliHelp()
	}
}
