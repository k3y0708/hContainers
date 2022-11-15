package main

import "fmt"

func cliHelp() {
	fmt.Println("Usage: hContainers <command> [arguments]")
	fmt.Println("Commands:")
	runnerHelp()
	containerHelp()
	fmt.Println("  help                                                                - Show this help message")
}

func cliContainersHelp() {
	fmt.Println("Usage: hContainers containers [arguments]")
	fmt.Println("Commands:")
	containerHelp()
}

func cliRunnerHelp() {
	fmt.Println("Usage: hContainers runner [arguments]")
	fmt.Println("Commands:")
	runnerHelp()
}

func runnerHelp() {
	fmt.Println("  runner list                                                         - List all available runners")
	fmt.Println("  runner create <runner-name>                                         - Create a new runner")
	fmt.Println("  runner delete <runner-name>                                         - Delete a runner")
}

func containerHelp() {
	fmt.Println("  containers list                                                      - List all available containers")
	fmt.Println("  containers create <runner-name> <container-name> <image>             - Create a new container")
	fmt.Println("  containers create <runner-name> <container-name> <image> <port:port> - Create a new container with port mapping")
	fmt.Println("  containers delete <name>                                             - Delete a container")
	fmt.Println("  containers start <name>                                              - Start a container")
	fmt.Println("  containers stop <name>                                               - Stop a container")
	fmt.Println("  containers restart <name>                                            - Restart a container")
	fmt.Println("  containers pause <name> <command>                                    - Pause a container")
	fmt.Println("  containers unpause <name>                                            - Unpause a container")
	fmt.Println("  containers exec <name> <command>                                     - Execute a command in a container (Non-interactive)")
	fmt.Println("  containers logs <name>                                               - Show logs of a container")
}
