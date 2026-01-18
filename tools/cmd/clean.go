package cmd

import (
	"fmt"
	"os"

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
	dirs := []string{"Core", "Drivers", "Middlewares", "STM32_BLE", "Library", "build"}
	files := []string{"Makefile"}

	fmt.Println("Cleaning project...")
	for _, d := range dirs {
		if err := os.RemoveAll(d); err != nil {
			fmt.Printf("Error removing %s: %v\n", d, err)
		} else {
			fmt.Printf("Removed %s\n", d)
		}
	}
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
