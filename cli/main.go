package main

import (
	"fmt"
	"os"
	"os/exec" // lets go run shell scripts
	"strings"
)

func main() {
	fmt.Println("SwitchConfigSim CLI Tool")
	fmt.Println()

	// Check if user provided command line args
	if len(os.Args) < 2 {
		showUsage()
		return
	}

	// command will be first arg after program name
	command := os.Args[1]

	switch command {
	case "version":
		handleVersionCommand()
	case "show":
		handleShowCommand()
	case "set":
		handleSetCommand()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run 'switchctl' without args to see available commands")
	}
}

// func to display usage info
func showUsage() {
	fmt.Println("Usage: switchctl <command> [options]")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  show                              - Display complete switch status")
	fmt.Println("  show system                       - Show system information only")
	fmt.Println("  show interfaces                   - Show interface status only")
	fmt.Println("  set hostname <name>               - Set switch hostname")
	fmt.Println("  set interface <name> <up/down>    - Set interface state")
	fmt.Println("  version                           - Show version information")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  switchctl show")
	fmt.Println("  switchctl show interfaces")
	fmt.Println("  switchctl set hostname core-switch-01")
	fmt.Println("  switchctl set interface eth0 up")
	fmt.Println("  switchctl set interface eth1 down")
}

// Function to handle version command
func handleVersionCommand() {
	fmt.Println("SwitchConfigSim v1.0.0")
	fmt.Println("Built with Go")
	fmt.Println("Simulates NVIDIA NVUE switch management")
}

// Function to handle show commands
func handleShowCommand() {
	// stores which part of status to show, either show system or interfaces
	var component string

	// Check if user typed a third word like system or interfaces
	if len(os.Args) >= 3 {
		component = os.Args[2]
	}

	// call the status script with appropriate component
	scriptPath := "./shell/get_status.sh"
	// stores the right script to run
	var cmd *exec.Cmd

	if component != "" {
		cmd = exec.Command(scriptPath, component)
		// runs the script with the component
	} else {
		cmd = exec.Command(scriptPath)
	}

	// Runs the script and captures what it prints
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing status script: %v\n", err)
		fmt.Printf("Script output: %s\n", string(output))
		return
	}

	// Display the script output
	fmt.Print(string(output))
}

// func to handle set commands
func handleSetCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Error: 'set' command requires additional arguments")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  switchctl set hostname <name>")
		fmt.Println("  switchctl set interface <name> <up|down>")
		return
	}
	// gets the third word after switchctl set
	// example: switchctl set hostname core-switch-01, third word is hostname
	subCommand := os.Args[2]

	switch subCommand {
	case "hostname":
		handleSetHostname()
	case "interface":
		handleSetInterface()
	default:
		fmt.Printf("Unknown set command: %s\n", subCommand)
		fmt.Println("Valid set commands: hostname, interface")
	}
}

// func to handle hostname setting
func handleSetHostname() {
	// checks if user typed a hostname
	if len(os.Args) < 4 {
		fmt.Println("Error: hostname command requires a hostname")
		fmt.Println("Usage: switchctl set hostname <name>")
		fmt.Println("Example: switchctl set hostname core-switch-01")
		return
	}
	// stores actual hostname
	hostname := os.Args[3]

	// call the hostname script
	scriptPath := "./shell/set_hostname.sh"
	cmd := exec.Command(scriptPath, hostname)

	// exec script and get output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error setting hostname: %v\n", err)
		fmt.Printf("Script output: %s\n", string(output))
		return
	}

	// Display success message
	fmt.Print(string(output))
}

// func to handle interface setting
func handleSetInterface() {
	// checks if user typed an interface name and state
	if len(os.Args) < 5 {
		fmt.Println("Error: interface command requires interface name and state")
		fmt.Println("Usage: switchctl set interface <name> <up|down>")
		fmt.Println("Example: switchctl set interface eth0 up")
		return
	}

	interfaceName := os.Args[3]
	interfaceState := os.Args[4]

	// validate state arg
	state := strings.ToLower(interfaceState)
	if state != "up" && state != "down" {
		fmt.Printf("Error: Invalid interface state '%s'\n", interfaceState)
		fmt.Println("Valid states: up, down")
		return
	}

	// call the interface script
	scriptPath := "./shell/set_interface.sh"
	cmd := exec.Command(scriptPath, interfaceName, state)

	// execute the script and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error setting interface: %v\n", err)
		fmt.Printf("Script output: %s\n", string(output))
		return
	}

	// Display success message
	fmt.Print(string(output))
}
