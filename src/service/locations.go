package service

import (
	"context"
	"fmt"
	"k3y0708/hContainers/global"
	"k3y0708/hContainers/util"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func LocationList() {
	locations := getAllLocations()
	fmt.Println("Available locations:")
	for _, location := range locations {
		if location.Name != "ash" {
			fmt.Printf("- %s\n", location.City)
		}
	}
}

func getLocationByName(locationName string) *hcloud.Location {
	locations := getAllLocations()

	for _, location := range locations {
		if location.City == locationName {
			return location
		}
	}
	return nil
}

func getAllLocations() []*hcloud.Location {
	locations, err := global.Client.Location.All(context.Background())
	util.CheckError(err, "Failed to get locations", -1)
	return locations
}
