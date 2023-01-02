package assembler

import "github.com/hetznercloud/hcloud-go/hcloud"

func ServerListToNames(servers []*hcloud.Server) []string {
	names := make([]string, len(servers))
	for i, server := range servers {
		names[i] = server.Name
	}
	return names
}
