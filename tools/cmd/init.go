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

	var preservableFiles = []string{
		"Core/Src/main.c",
		"Core/Inc/main.h",
		"Core/Inc/app_conf.h",
		"Core/Src/stm32wb0x_hal_msp.c",
		"Core/Src/app_entry.c",
		"STM32_BLE/App/app_ble.c",
		"STM32_BLE/App/ble_conf.h",
	}

	savedCodes := make(map[string]internal.UserCodeMap)

	fmt.Println("Preserving user code...")
	for _, f := range preservableFiles {
		codes, err := internal.ExtractUserCode(f)
		if err == nil && len(codes) > 0 {
			savedCodes[f] = codes
			fmt.Printf("  Saved %d blocks from %s\n", len(codes), f)
		}
	}

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

	// Restore User Code
	fmt.Println("Restoring user code...")
	for f, codes := range savedCodes {
		if err := internal.RestoreUserCodeInFile(f, codes); err != nil {
			fmt.Printf("  Error restoring %s: %v\n", f, err)
		} else {
			fmt.Printf("  Restored %s\n", f)
		}
	}

		// 5. Generate Makefile
		fmt.Println("Generating Makefile...")
		if err := os.WriteFile("Makefile", []byte(templates.Makefile), 0644); err != nil {
			fmt.Printf("Error generating Makefile: %v\n", err)
		}
	
		// 6. Apply Project Configuration
		if config, err := internal.LoadConfig(); err == nil {
			fmt.Println("Applying project configuration...")
			internal.ApplyLPMSetting(config.DisableLPM)
		}
		
		fmt.Println("Initialization complete.")
	}
