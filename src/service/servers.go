package service

import (
	"context"

	"github.com/hContainers/hContainers/global"
	"github.com/hContainers/hContainers/util"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func GetAllServers() []*hcloud.Server {
	servers, _, err := global.Client.Server.List(context.Background(), hcloud.ServerListOpts{})
	util.CheckError(err, "Failed to get servers", -1)
	return servers
}

func GetAllServerNames() []string {
	var serverNames []string
	for _, server := range GetAllServers() {
		serverNames = append(serverNames, server.Name)
	}
	return serverNames
}

func checkIfServerExists(serverName string) bool {
	servers := GetAllServers()

	for _, server := range servers {
		if server.Name == serverName {
			return true
		}
	}
	return false
}
