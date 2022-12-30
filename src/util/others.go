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

/**
 * Returns an int representing how different two strings are
 * The lower the number, the more similar the strings are
 */
func StringSimilarity(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	var cost int
	if a[len(a)-1] == b[len(b)-1] {
		cost = 0
	} else {
		cost = 1
	}

	res := min(
		StringSimilarity(a[:len(a)-1], b)+1,
		StringSimilarity(a, b[:len(b)-1])+1,
		StringSimilarity(a[:len(a)-1], b[:len(b)-1])+cost,
	)

	return res
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	} else if b < c {
		return b
	} else {
		return c
	}
}

func FindNearestCommand(command string, possibleCommands []string) string {
	nearestCommand := ""
	nearestCommandDistance := 4

	for _, possibleCommand := range possibleCommands {
		distance := StringSimilarity(command, possibleCommand)
		if distance < nearestCommandDistance {
			nearestCommand = possibleCommand
			nearestCommandDistance = distance
		}
	}

	return nearestCommand
}
