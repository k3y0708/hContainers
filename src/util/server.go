package util

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hContainers/hContainers/global"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func GetAllServers() []*hcloud.Server {
	servers, _, err := global.Client.Server.List(context.Background(), hcloud.ServerListOpts{})
	CheckError(err, "Failed to get servers", -1)
	return servers
}

func GetAllServerNames() []string {
	var serverNames []string
	for _, server := range GetAllServers() {
		serverNames = append(serverNames, server.Name)
	}
	return serverNames
}

func GetServerByName(serverName string) *hcloud.Server {
	servers, _, err := global.Client.Server.List(context.Background(), hcloud.ServerListOpts{})
	CheckError(err, "Failed to get servers", -1)
	for _, server := range servers {
		if server.Name == serverName {
			return server
		}
	}
	return nil
}

func GetIpByName(serverName string) string {
	return GetServerByName(serverName).PublicNet.IPv4.IP.String()
}

func CheckIfServerExists(serverName string) bool {
	servers := GetAllServerNames()

	for _, server := range servers {
		if server == serverName {
			return true
		}
	}
	return false
}

func CheckIfRunnerNameIsValid(runnerName string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9-.]*$`).MatchString(runnerName)
}

func WaitForServerToBeOnline(serverName string) error {
	status := hcloud.ServerStatusUnknown // Waiting for server to be down :D
	sshPossible := false
	ip := GetIpByName(serverName)
	for i := 0; i < 120; i++ {
		if status == hcloud.ServerStatusRunning {
			// Wait until ssh connection is possible
			if CheckIfSshConnectionIsPossible(ip) {
				sshPossible = true
				break
			}
		}
		fmt.Printf("Waiting for server %s to be online...\n", serverName)
		time.Sleep(5 * time.Second)
		status = GetServerByName(serverName).Status
	}
	if sshPossible {
		return nil
	} else {
		return fmt.Errorf("server %s is not online", serverName)
	}
}

func CheckIfSshConnectionIsPossible(serverIP string) bool {
	stdout, err := RemoteRun(serverIP, "echo 'SSH connection is possible!'", 0)
	if err != nil {
		return false
	}
	if stdout == "SSH connection is possible!\n" {
		return true
	}
	fmt.Println(stdout)
	return false
}
