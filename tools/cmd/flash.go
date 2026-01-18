package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/spf13/cobra"
)

var flashCmd = &cobra.Command{
	Use:   "flash",
	Short: "Build and Flash the project using CMake and Ninja",
	Run: func(cmd *cobra.Command, args []string) {
		runFlash()
	},
}

func init() {
	rootCmd.AddCommand(flashCmd)
}

func runFlash() {
	// 1. Configure if needed (check for build/build.ninja)
	if _, err := os.Stat("build/build.ninja"); os.IsNotExist(err) {
		fmt.Println("Configuring CMake with Ninja...")
		runCommand("cmake", "-DCMAKE_TOOLCHAIN_FILE=cmake/arm-none-eabi-gcc.cmake", "-G", "Ninja", "-B", "build", "-S", ".")
	}

	// 2. Build and Flash
	fmt.Println("Building and Flashing...")
	runCommand("cmake", "--build", "build", "--target", "flash")
}

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)

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
			fmt.Println(scanner.Text())
		}
	}

	go process(stdout)
	go process(stderr)

	wg.Wait()
	if err := cmd.Wait(); err != nil {
		// Exit with error code if the command failed
		os.Exit(1)
	}
}