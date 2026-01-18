package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format C/C++ and Go source files",
	Run: func(cmd *cobra.Command, args []string) {
		runFormat()
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)
}

func runFormat() {
	// 1. Format C/C++ files
	fmt.Println("Formatting C/C++ files...")
	// Focus on user code directories to avoid reformatting vendor libraries (Drivers, Middlewares)
	cDirs := []string{"Core", "Library", "STM32_BLE"}
	var cFiles []string

	for _, dir := range cDirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				ext := strings.ToLower(filepath.Ext(path))
				if ext == ".c" || ext == ".h" || ext == ".cpp" || ext == ".hpp" {
					cFiles = append(cFiles, path)
				}
			}
			return nil
		})
		if err != nil {
			// Just warn if a directory doesn't exist (e.g. strict subset of folders)
			// But don't crash
			fmt.Printf("Warning: Could not check %s: %v\n", dir, err)
		}
	}

	if len(cFiles) > 0 {
		// Pass -i to modify files in place
		// We pass all files at once. If there are too many, we might hit arg limit,
		// but for a MCU project it's likely fine.
		// If needed we could batch them.
		args := append([]string{"-i"}, cFiles...)
		cmd := exec.Command("clang-format", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error running clang-format: %v\n", err)
		} else {
			fmt.Printf("Formatted %d C/C++ files.\n", len(cFiles))
		}
	} else {
		fmt.Println("No C/C++ files found in user directories.")
	}

	// 2. Format Go files
	fmt.Println("Formatting Go files...")
	if _, err := os.Stat("tools"); err == nil {
		cmd := exec.Command("go", "fmt", "./...")
		cmd.Dir = "tools"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error running go fmt: %v\n", err)
		} else {
			fmt.Println("Go files formatted.")
		}
	} else {
		fmt.Println("Tools directory not found, skipping Go formatting.")
	}
}
