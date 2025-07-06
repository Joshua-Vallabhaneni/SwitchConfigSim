package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("SwitchConfigSim CLI Tool")
	fmt.Println()
	// Check if user provided command line args
	// provide instructions if no args provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: switchctl <command>")
		fmt.Println()
		fmt.Println("Available commands:")
		fmt.Println("  show     - Display switch")
		fmt.Println("  set      - Configure switch settings")
		fmt.Println("  version  - Show version info")
		fmt.Println()
		return
	}
	// command will be first arg after program name
	command := os.Args[1]

	switch command {
	case "version":
		fmt.Println("SwitchConfigSim v1.0.0")
		fmt.Println("Built with Go")
	case "show":
		handleShowCommand()
	case "set":
		handleSetCommand()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run 'switchctl' w/o args to see available commands")
	}
}

func handleShowCommand() {
	fmt.Println("Show - Read and Display switch state")
	fmt.Println()
	fmt.Println("Current Switch Configuration:")
	fmt.Println("Hostname: switch1")
	fmt.Println("Interfaces: eth0, eth1, eth2, ethx")
	// each switch has multiple ports, each port can be connected to something else
	// eth0 = first ethernet interface
	fmt.Println("Status: Operational")
	// May also be down, failed, in maintainece
}

func handleSetCommand() {
	fmt.Println("Set - Config switch settings")
	fmt.Println()
	fmt.Println("Hostname: switch1 -> datacenter-switch-1")
	fmt.Println("Interface eth1: down -> up")
	// example of port being enabled
	fmt.Println("Interface eth2: 1G -> 10G")
	// example of a speed upgarde 
	fmt.Println("Status: Opertional")
}
