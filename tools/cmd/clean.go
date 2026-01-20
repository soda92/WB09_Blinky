package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean project directories",
	Run: func(cmd *cobra.Command, args []string) {
		runClean()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func runClean() {
	dirs := []string{"Drivers", "Middlewares", "Library", "build"}
	files := []string{"Makefile", "compile_commands.json"}

	fmt.Println("Cleaning project (preserving Core/ and STM32_BLE/)...")
	for _, d := range dirs {
		if err := os.RemoveAll(d); err != nil {
			fmt.Printf("Error removing %s: %v\n", d, err)
		} else {
			fmt.Printf("Removed %s\n", d)
		}
	}
	
	// Also clean build artifacts in root if any
	buildArtifacts, _ := filepath.Glob("*.elf")
	buildArtifacts2, _ := filepath.Glob("*.bin")
	buildArtifacts3, _ := filepath.Glob("*.hex")
	buildArtifacts4, _ := filepath.Glob("*.map")
	files = append(files, buildArtifacts...)
	files = append(files, buildArtifacts2...)
	files = append(files, buildArtifacts3...)
	files = append(files, buildArtifacts4...)

	for _, f := range files {
		if err := os.Remove(f); err != nil {
			if !os.IsNotExist(err) {
				fmt.Printf("Error removing %s: %v\n", f, err)
			}
		} else {
			fmt.Printf("Removed %s\n", f)
		}
	}
	fmt.Println("Clean complete.")
}
