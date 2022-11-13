package util

import (
	"fmt"
	"os"
)

func CheckError(err error, message string, exitCode int) {
	if err != nil {
		fmt.Println(message)
		fmt.Println(err)
		os.Exit(exitCode)
	}
}
