package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func RunFlash() {
	cmd := exec.Command("make", "flash")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error getting stdout pipe: %v\n", err)
		os.Exit(1)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error getting stderr pipe: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	process := func(r io.Reader) {
		defer wg.Done()
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			cleanLine := strings.TrimSpace(line)

			// Heuristic to detect compiler commands
			isCompilerCmd := strings.HasPrefix(cleanLine, "arm-none-eabi-gcc") ||
				strings.HasPrefix(cleanLine, "arm-none-eabi-objcopy") ||
				strings.HasPrefix(cleanLine, "arm-none-eabi-size") ||
				strings.HasPrefix(cleanLine, "mkdir") || 
				strings.HasPrefix(cleanLine, "rm -")

			// Always show errors and warnings
			isErrorOrWarning := strings.Contains(strings.ToLower(line), "error:") ||
				strings.Contains(strings.ToLower(line), "warning:")

			// Always show our custom echo messages and Programmer output headers
			isEssential := strings.Contains(line, "Flashing via") ||
				strings.Contains(line, "Done!") ||
				strings.Contains(line, "ST-LINK SN") ||
				strings.Contains(line, "Device name") ||
				strings.Contains(line, "Flash size") ||
				strings.Contains(line, "Download verified successfully")

			if isErrorOrWarning || isEssential || !isCompilerCmd {
				fmt.Println(line)
			}
		}
	}

	go process(stdout)
	go process(stderr)

	wg.Wait()
	if err := cmd.Wait(); err != nil {
		// Just exit with the same code if possible, or 1
		os.Exit(1)
	}
}
