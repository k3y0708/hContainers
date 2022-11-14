package util

import "fmt"

func ShowHelp() {
	fmt.Println("Usage: hContainers <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  runner list                                                         - List all available runners")
	fmt.Println("  runner create <runner-name>                                         - Create a new runner")
	fmt.Println("  runner delete <runner-name>                                         - Delete a runner")
	fmt.Println("  container list                                                      - List all available containers")
	fmt.Println("  container create <runner-name> <container-name> <image>             - Create a new container")
	fmt.Println("  container create <runner-name> <container-name> <image> <port:port> - Create a new container with port mapping")
	fmt.Println("  container delete <name>                                             - Delete a container")
	fmt.Println("  container start <name>                                              - Start a container")
	fmt.Println("  container stop <name>                                               - Stop a container")
	fmt.Println("  container restart <name>                                            - Restart a container")
	fmt.Println("  container pause <name> <command>                                    - Pause a container")
	fmt.Println("  container unpause <name>                                            - Unpause a container")
	fmt.Println("  container exec <name> <command>                                     - Execute a command in a container (Non-interactive)")
	fmt.Println("  container logs <name>                                               - Show logs of a container")
	fmt.Println("  help                                                                - Show this help message")
}
