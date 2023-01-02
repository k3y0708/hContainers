package service

import (
	"context"

	"github.com/hContainers/hContainers/global"
	"github.com/hContainers/hContainers/util"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func getLocationByName(locationName string) *hcloud.Location {
	locations := GetAllLocations()

	for _, location := range locations {
		if location.City == locationName {
			return location
		}
	}
	return nil
}

func GetAllLocations() []*hcloud.Location {
	locations, err := global.Client.Location.All(context.Background())
	util.CheckError(err, "Failed to get locations", -1)
	return locations
}
