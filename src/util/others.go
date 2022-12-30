package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/hContainers/hContainers/colors"
)

func CheckError(err error, message string, exitCode int) {
	if err != nil {
		fmt.Println(message)
		fmt.Println(err)
		os.Exit(exitCode)
	}
}

func CheckLength(slice []string, length int, message string, exitCode int) {
	if len(slice) < length {
		fmt.Println(message)
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

func Copy(src string, dst string) {
	data, err := os.ReadFile(src)
	CheckError(err, "Failed to read file", 1)
	err = os.WriteFile(dst, data, 0755)
	CheckError(err, "Failed to write file", 1)
}
