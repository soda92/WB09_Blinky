package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wb09_tool",
	Short: "A tool for STM32WB09 project management",
	Long:  `A comprehensive tool to initialize, manage dependencies, and flash STM32WB09 projects.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
