package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"wb09_tool/internal"
	"wb09_tool/templates"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the project from template",
	Run: func(cmd *cobra.Command, args []string) {
		runInit()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit() {
	templatePath := internal.DefaultTemplatePath
	fmt.Printf("Initializing project from template: %s\n", templatePath)

	// 1. Copy Application Code (Core, STM32_BLE) from Template
	fmt.Println("Copying Application Code...")
	if err := internal.CopyDir(filepath.Join(templatePath, "Core"), "Core"); err != nil {
		fmt.Printf("Error copying Core: %v\n", err)
	}
	if err := internal.CopyDir(filepath.Join(templatePath, "STM32_BLE"), "STM32_BLE"); err != nil {
		fmt.Printf("Error copying STM32_BLE: %v\n", err)
	}

	// 2. Copy SDK Components (Drivers, Middlewares)
	fmt.Println("Copying SDK Components...")
	if err := internal.CopyDir(filepath.Join(internal.SDKPath, "Drivers"), "Drivers"); err != nil {
		fmt.Printf("Error copying Drivers: %v\n", err)
	}
	
	// Create Middlewares/ST structure before copying
	os.MkdirAll("Middlewares/ST", 0755)
	if err := internal.CopyDir(filepath.Join(internal.SDKPath, "Middlewares/ST/STM32_BLE"), "Middlewares/ST/STM32_BLE"); err != nil {
		fmt.Printf("Error copying Middlewares: %v\n", err)
	}

	// 3. Remove Templates from Drivers to avoid build errors
	fmt.Println("Cleaning up Drivers templates...")
	driversSrc := filepath.Join(internal.DriversDir, "STM32WB0x_HAL_Driver/Src")
	files, _ := filepath.Glob(filepath.Join(driversSrc, "*_template.c"))
	for _, f := range files {
		os.Remove(f)
	}

	// 4. Create Library Directories
	os.MkdirAll(internal.LibSrcDir, 0755)
	os.MkdirAll(internal.LibIncDir, 0755)

	// 5. Generate Makefile
	fmt.Println("Generating Makefile...")
	if err := os.WriteFile("Makefile", []byte(templates.Makefile), 0644); err != nil {
        fmt.Printf("Error generating Makefile: %v\n", err)
    }
	
	fmt.Println("Initialization complete.")
}
