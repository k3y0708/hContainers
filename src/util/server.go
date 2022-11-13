package util

import (
	"context"
	"regexp"

	"k3y0708/hContainers/global"

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

func GetAmountOfRunners() int {
	servers, _, err := global.Client.Server.List(context.Background(), hcloud.ServerListOpts{ListOpts: hcloud.ListOpts{LabelSelector: "runner"}})
	CheckError(err, "Failed to get servers", -1)
	return len(servers)
}
