package util

import (
	"fmt"
	"k3y0708/hContainers/colors"
	"os"
	"strings"
)

func CheckError(err error, message string, exitCode int) {
	if err != nil {
		fmt.Println(message)
		fmt.Println(err)
		os.Exit(exitCode)
	}
}

func StatusToColor(status string) string {
	switch strings.ToLower(status) {
	case "up":
		return colors.GREEN
	case "exited":
		return colors.RED
	case "paused":
		return colors.YELLOW
	default:
		return colors.WHITE
	}
}
