package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor [port] [baud]",
	Short: "Monitor Serial Output for 10 seconds using picocom",
	Run: func(cmd *cobra.Command, args []string) {
		port := "/dev/ttyACM0"
		if len(args) > 0 {
			port = args[0]
		}
		baud := "115200"
		if len(args) > 1 {
			baud = args[1]
		}

		fmt.Printf("Monitoring %s at %s for 10 seconds...\n", port, baud)

		// Check if port exists
		if _, err := os.Stat(port); os.IsNotExist(err) {
			fmt.Printf("Error: Serial port %s not found.\n", port)
			// Try finding standard STM32 port
			fmt.Println("Searching for available ttyACM ports...")
			matches, _ := filepath.Glob("/dev/ttyACM*")
			if len(matches) > 0 {
				fmt.Printf("Found: %v. Using %s\n", matches, matches[0])
				port = matches[0]
			} else {
				return
			}
		}

		c := exec.Command("timeout", "10s", "picocom", "-b", baud, "-q", port)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		
		err := c.Run()
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 124 {
				fmt.Println("\n--- Monitoring finished (10s limit) ---")
				return
			}
		}
		
		if err != nil {
			fmt.Printf("Error running monitor: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
