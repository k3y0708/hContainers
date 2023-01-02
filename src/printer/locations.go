package printer

import "fmt"

func LocationsList(locations []string) {
	for _, location := range locations {
		fmt.Println(location)
	}
}
