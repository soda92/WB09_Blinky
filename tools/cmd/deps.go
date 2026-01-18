package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"wb09_tool/internal"
)

var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Scan sources and copy missing dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		runDeps()
	},
}

func init() {
	rootCmd.AddCommand(depsCmd)
}

var processedFiles = make(map[string]bool)

// Directories to check before deciding a file is missing
var includePaths = []string{
	"Core/Inc",
	"Library/Inc",
	"STM32_BLE/App",
	"STM32_BLE/Target",
	"Drivers/STM32WB0x_HAL_Driver/Inc",
	"Drivers/STM32WB0x_HAL_Driver/Inc/Legacy",
	"Drivers/BSP/STM32WB0x-nucleo",
	"Drivers/CMSIS/Device/ST/STM32WB0X/Include",
	"Drivers/CMSIS/Include",
	"Middlewares/ST/STM32_BLE",
	"Middlewares/ST/STM32_BLE/stack/include",
	"Middlewares/ST/STM32_BLE/evt_handler/inc",
	"Middlewares/ST/STM32_BLE/cryptolib/Inc",
}

var sourcePaths = []string{
	"Core/Src",
	"Library/Src",
	"STM32_BLE/App",
	"STM32_BLE/Target",
	"Drivers/STM32WB0x_HAL_Driver/Src",
	"Drivers/BSP/STM32WB0x-nucleo",
	"Middlewares/ST/STM32_BLE/evt_handler/src",
	"Middlewares/ST/STM32_BLE/stack/config",
}

var knownSymbols = map[string]string{
	"blue_unit_conversion": "blue_unit_conversion.s",
	"CPUcontextSave":       "cpu_context_switch.s",
	"APP_DEBUG_SIGNAL_SET": "app_debug.c",
	"RT_DEBUG_GPIO_Init":   "app_debug.c",
    "HOST_TO_LE_16":        "ble_types.h", // Just in case
}

func runDeps() {
	os.MkdirAll(internal.LibIncDir, 0755)
	os.MkdirAll(internal.LibSrcDir, 0755)
	
	// Seed with existing source files in all source paths
	for _, p := range sourcePaths {
		files, _ := filepath.Glob(filepath.Join(p, "*.c"))
		for _, f := range files {
			scanFile(f)
		}
	}
}

func scanFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Check for Includes
		if strings.HasPrefix(line, "#include") {
			parts := strings.Split(line, "\"")
			if len(parts) > 1 {
				header := parts[1]
				processDependency(header)
			}
		}

		// Check for Known Symbols
		for symbol, filename := range knownSymbols {
			if strings.Contains(line, symbol) {
				processDependency(filename)
			}
		}
	}
}

func processDependency(filename string) {
	if processedFiles[filename] {
		return
	}
	processedFiles[filename] = true

	isHeader := strings.HasSuffix(filename, ".h")
	var targetDir string
	var checkPaths []string

	if isHeader {
		targetDir = internal.LibIncDir
		checkPaths = includePaths
	} else {
		targetDir = internal.LibSrcDir
		checkPaths = sourcePaths
	}

	// Check if exists in paths
	existingPath := internal.FindInPaths(filename, checkPaths)
	if existingPath != "" {
		// Scan it for recursive includes/symbols
		scanFile(existingPath)
		return
	}
	
	// Not found locally, search in Repo
	path := internal.FindFileInRepo(internal.SDKPath, filename)
	if path != "" {
		fmt.Printf("Found missing dependency: %s -> Copying to %s\n", filename, targetDir)
		internal.CopyFile(path, filepath.Join(targetDir, filename))
		
		// Scan the new file
		scanFile(filepath.Join(targetDir, filename))

		if isHeader {
			// Check for associated .c file
			cFile := strings.Replace(filename, ".h", ".c", 1)
			// Check if .c file already exists or processed
			if internal.FindInPaths(cFile, sourcePaths) == "" && !processedFiles[cFile] {
				// Search and copy associated .c file
				cPath := filepath.Join(filepath.Dir(path), cFile)
				if !internal.FileExists(cPath) {
					cPath = internal.FindFileInRepo(internal.SDKPath, cFile)
				}
				if cPath != "" {
					fmt.Printf("Found associated source: %s -> Copying to %s\n", cFile, internal.LibSrcDir)
					internal.CopyFile(cPath, filepath.Join(internal.LibSrcDir, cFile))
					processedFiles[cFile] = true
					scanFile(filepath.Join(internal.LibSrcDir, cFile))
				}
			}
		}
	}
}
