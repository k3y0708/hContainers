package assembler

import "github.com/hetznercloud/hcloud-go/hcloud"

func LocationsToStrings(locations []*hcloud.Location) []string {
	var strings []string
	for _, location := range locations {
		strings = append(strings, location.Name)
	}
	return strings
}
