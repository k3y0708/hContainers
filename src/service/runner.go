package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hContainers/hContainers/global"
	"github.com/hContainers/hContainers/types"
	"github.com/hContainers/hContainers/util"

	"github.com/hetznercloud/hcloud-go/hcloud"

	_ "embed"
)

//go:embed cloudinit.yml
var cloudinit string

func RunnerRestart(name string) {
	_, _, err := global.Client.Server.Reboot(context.Background(), util.GetServerByName(name))
	util.CheckError(err, "Failed to restart runner", 1)
	// Wait for server to be online
	err = util.WaitForServerToBeOnline(name)
	util.CheckError(err, "Failed to wait for runner to be online", 1)
	fmt.Println("Restarting all containers on runner")
	runnerStartAllContainers(name)
}

func runnerStartAllContainers(runnerName string) {
	containers := getContainersFromRunner(util.GetServerByName(runnerName))
	for _, container := range containers {
		ContainerStart(container.Name)
	}
}

func RunnerDelete(name string) {
	_, _, err := global.Client.Server.DeleteWithResult(context.Background(), util.GetServerByName(name))
	util.CheckError(err, "Failed to delete runner", 1)
}

func RunnerCreate(name string, flags types.FlagsServer) {
	if util.CheckIfServerExists(name) {
		fmt.Println("Runner already exists")
		os.Exit(1)
	}
	location := getLocationByName(flags.Location)
	create_options := hcloud.ServerCreateOpts{
		Name:       name,
		ServerType: &hcloud.ServerType{Name: flags.Sku},
		Image:      &hcloud.Image{Name: "ubuntu-22.04"},
		Location:   location,
		UserData:   strings.Replace(cloudinit, "{{{PUBLIC_KEY}}}", global.PublicKey, 1),
		PublicNet:  &hcloud.ServerCreatePublicNet{EnableIPv4: true, EnableIPv6: !flags.DisableIPv6},
		Labels:     map[string]string{"runner": "true"},
	}
	_, _, err := global.Client.Server.Create(context.Background(), create_options)
	util.CheckError(err, "Failed to create runner", 1)
}
