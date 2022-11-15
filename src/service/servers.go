package service

import (
	"context"
	"k3y0708/hContainers/global"
	"k3y0708/hContainers/util"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func GetAllServers() []*hcloud.Server {
	servers, _, err := global.Client.Server.List(context.Background(), hcloud.ServerListOpts{})
	util.CheckError(err, "Failed to get servers", -1)
	return servers
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
