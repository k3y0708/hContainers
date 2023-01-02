package printer

import (
	"fmt"
)

func RunnerList(runnernames []string) {
	fmt.Printf("Available runners (Amount: %d):\n", len(runnernames))
	for _, runner := range runnernames {
		fmt.Printf("- %s\n", runner)
	}
}
