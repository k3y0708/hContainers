package main

import (
	"fmt"
	"strings"

	"github.com/hContainers/hContainers/types"
)

/*
Prints the help message for the help command
*/
func cliHelp() {
	fmt.Println("Usage: hContainers <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println(formatter(helpHelp()))
}

/*
Prints the help message for the containers command
*/
func cliContainersHelp() {
	fmt.Println("Usage: hContainers containers [arguments]")
	fmt.Println("Commands:")
	fmt.Println(formatter(containerHelp()))
}

/*
Prints the help message for the runners command
*/
func cliRunnerHelp() {
	fmt.Println("Usage: hContainers runner [arguments]")
	fmt.Println("Commands:")
	fmt.Println(formatter(runnerHelp()))
}

/*
Prints the help message for the locations command
*/
func cliLocationsHelp() {
	fmt.Println("Usage: hContainers locations [arguments]")
	fmt.Println("Commands:")
	fmt.Println(formatter(locationsHelp()))
}

/*
Returns the help rows for the help command

@return The help rows
*/
func helpHelp() []types.HelpRow {
	return []types.HelpRow{
		{Command: "help", Description: "Show this help message"},
		{Command: "version", Description: "Show version"},
		{Command: "runners", Description: "Manage runners"},
		{Command: "containers", Description: "Manage containers"},
		{Command: "locations", Description: "Manage locations"},
	}
}

/*
Returns the help rows for the runners command

@return The help rows
*/
func runnerHelp() []types.HelpRow {
	return []types.HelpRow{
		{Command: "runners list", Description: "List all available runners"},
		{Command: "runners create <runner-name>", Description: "Create a new runner"},
		{Command: "runners delete <runner-name>", Description: "Delete a runner"},
	}
}

/*
Returns the help rows for the containers command

@return The help rows
*/
func containerHelp() []types.HelpRow {
	return []types.HelpRow{
		{Command: "containers list", Description: "List all available containers"},
		{Command: "containers create <runner-name> <container-name> <image>", Description: "Create a new container"},
		{Command: "containers delete <name>", Description: "Delete a container"},
		{Command: "containers start <name>", Description: "Start a container"},
		{Command: "containers stop <name>", Description: "Stop a container"},
		{Command: "containers restart <name>", Description: "Restart a container"},
		{Command: "containers pause <name> <command>", Description: "Pause a container"},
		{Command: "containers unpause <name>", Description: "Unpause a container"},
		{Command: "containers exec <name> <command>", Description: "Execute a command in a container (Non-interactive)"},
		{Command: "containers logs <name>", Description: "Show logs of a container"},
	}
}

/*
Returns the help rows for the locations command

@return The help rows
*/
func locationsHelp() []types.HelpRow {
	return []types.HelpRow{
		{
			Command:     "locations list",
			Description: "List all available locations",
		},
	}
}

/*
Formats the help rows into a string

@param rows The rows to format

@return The formatted multi-line string
*/
func formatter(rows ...[]types.HelpRow) string {
	allRows := []types.HelpRow{}
	for _, row := range rows {
		allRows = append(allRows, row...)
	}

	out := ""

	longestCommand := 0
	for _, row := range allRows {
		if len(row.Command) > longestCommand {
			longestCommand = len(row.Command)
		}
	}

	for _, row := range allRows {
		out += "  " + row.Command + " " + strings.Repeat(" ", longestCommand-len(row.Command)) + " - " + row.Description + "\n"
	}
	return out
}
