package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hContainers/hContainers/util"

	"github.com/hContainers/hContainers/global"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

func main() {
	args := os.Args[1:]

	hcloudToken, present := os.LookupEnv(global.ENV_HETZNER_TOKEN)
	if !present {
		fmt.Println(global.ENV_HETZNER_TOKEN + " not set")
		os.Exit(1)
	}
	global.Client = hcloud.NewClient(hcloud.WithToken(hcloudToken))

	sshKeyPath, present := os.LookupEnv(global.ENV_SSH_KEY_PATH)
	if !present {
		fmt.Println(global.ENV_SSH_KEY_PATH + " not set")
		os.Exit(1)
	}
	sshKeyPath = strings.Replace(sshKeyPath, "~", os.Getenv("HOME"), 1)
	publicKeyBytes, err := os.ReadFile(sshKeyPath + ".pub")
	util.CheckError(err, "Failed to read public key", 1)
	global.PublicKey = string(publicKeyBytes)
	privateKeyBytes, err := os.ReadFile(sshKeyPath)
	util.CheckError(err, "Failed to read private key", 1)
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
	case "locations":
		cliLocations(args[1:])
	case "containers":
		cliContainers(args[1:])
	case "container":
		fmt.Println("Invalid command\nDid you mean 'containers'?")
	case "version":
		fmt.Println(global.Version)
	default:
		fmt.Println("Invalid command")
		cliHelp()
	}
}
