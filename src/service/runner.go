package service

import (
	"context"
	"fmt"
	"k3y0708/hContainers/global"
	"k3y0708/hContainers/util"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"

	_ "embed"
)

//go:embed cloudinit.yml
var cloudinit string

func RunnerRestart(name string) {
	fmt.Println("Restarting runner...")
	_, _, err := global.Client.Server.Reboot(context.Background(), util.GetServerByName(name))
	util.CheckError(err, "Failed to restart runner", 1)
	fmt.Println("Runner restarted")
}

func RunnerDelete(name string) {
	fmt.Println("Deleting runner...")
	_, _, err := global.Client.Server.DeleteWithResult(context.Background(), util.GetServerByName(name))
	util.CheckError(err, "Failed to delete runner", 1)
	fmt.Println("Runner deleted")
}

func RunnerCreate(name string) {
	fmt.Println("Creating runner...")
	create_options := hcloud.ServerCreateOpts{
		Name:       name,
		ServerType: &hcloud.ServerType{Name: "cx11"},
		Image:      &hcloud.Image{Name: "ubuntu-22.04"},
		Location:   &hcloud.Location{City: "fsn1"},
		UserData:   strings.Replace(cloudinit, "{{{PUBLIC_KEY}}}", global.PublicKey, 1),
		PublicNet:  &hcloud.ServerCreatePublicNet{EnableIPv4: true, EnableIPv6: true},
		Labels:     map[string]string{"runner": "true"},
	}
	_, _, err := global.Client.Server.Create(context.Background(), create_options)
	util.CheckError(err, "Failed to create runner", 1)
	fmt.Println("Runner created")
}
