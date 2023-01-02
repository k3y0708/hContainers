package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/hContainers/hContainers/colors"
	"github.com/hContainers/hContainers/global"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

func CheckError(err error, message string, exitCode int) {
	if err != nil {
		fmt.Println(message)
		fmt.Println(err)
		os.Exit(exitCode)
	}
}

func CheckLength(slice []string, length int, message string, exitCode int) {
	if len(slice) < length {
		fmt.Println(message)
		os.Exit(exitCode)
	}
}

func StatusToColor(status string) string {
	switch strings.ToLower(status) {
	case "up":
		return colors.GREEN
	case "exited":
		return colors.RED
	case "paused":
		return colors.YELLOW
	default:
		return colors.WHITE
	}
}

func Copy(src string, dst string) {
	data, err := os.ReadFile(src)
	CheckError(err, "Failed to read file", 1)
	err = os.WriteFile(dst, data, 0755)
	CheckError(err, "Failed to write file", 1)
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	} else if b < c {
		return b
	} else {
		return c
	}
}

func FindNearestCommand(command string, possibleCommands []string) string {
	nearestCommand := ""
	nearestCommandDistance := 4

	for _, possibleCommand := range possibleCommands {
		distance := stringSimilarity(command, possibleCommand)
		if distance < nearestCommandDistance {
			nearestCommand = possibleCommand
			nearestCommandDistance = distance
		}
	}

	return nearestCommand
}

func GetEnv() {
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
	CheckError(err, "Failed to read public key", 1)
	global.PublicKey = string(publicKeyBytes)
	privateKeyBytes, err := os.ReadFile(sshKeyPath)
	CheckError(err, "Failed to read private key", 1)
	global.PrivateKey = string(privateKeyBytes)
}
