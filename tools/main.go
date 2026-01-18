package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "setup":
		RunSetup()
	case "deps":
		RunDeps()
	case "flash":
		RunFlash()
	case "all":
		RunSetup()
		RunDeps()
		RunFlash()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run tools/*.go <command>")
	fmt.Println("Commands:")
	fmt.Println("  setup  - Copy Drivers and Middlewares from SDK")
	fmt.Println("  deps   - Scan sources and copy missing headers/sources")
	fmt.Println("  flash  - Build (make qflash logic) and Flash")
	fmt.Println("  all    - Run setup, deps, and flash in sequence")
}

