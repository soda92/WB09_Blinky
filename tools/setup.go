package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func RunSetup() {
	fmt.Println("Setting up project environment...")

	// Verify Repo Path
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fmt.Printf("Error: STM32Cube Repository not found at %s\n", repoPath)
		return
	}

	// Copy Drivers
	driversSrc := filepath.Join(repoPath, "Drivers")
	driversDst := "Drivers"
	if _, err := os.Stat(driversDst); os.IsNotExist(err) {
		fmt.Printf("Copying Drivers from %s...\n", driversSrc)
		cmd := exec.Command("cp", "-r", driversSrc, driversDst)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error copying Drivers: %v\n", err)
		}
	} else {
		fmt.Println("Drivers directory already exists. Skipping.")
	}

	// Copy Middlewares (specifically STM32_BLE)
	mwSrc := filepath.Join(repoPath, "Middlewares/ST/STM32_BLE")
	mwDst := "Middlewares/ST/STM32_BLE"
	if _, err := os.Stat(mwDst); os.IsNotExist(err) {
		fmt.Printf("Copying STM32_BLE Middleware from %s...\n", mwSrc)
		os.MkdirAll("Middlewares/ST", 0755)
		cmd := exec.Command("cp", "-r", mwSrc, mwDst)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error copying Middleware: %v\n", err)
		}
	} else {
		fmt.Println("Middlewares/ST/STM32_BLE already exists. Skipping.")
	}
	
	fmt.Println("Setup complete.")
}
