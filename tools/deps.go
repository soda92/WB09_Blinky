package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const repoPath = "/home/soda/STM32Cube/Repository/STM32Cube_FW_WB0_V1.4.0"
const destInc = "Library/Inc"
const destSrc = "Library/Src"

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
}

var processedFiles = make(map[string]bool)
var foundSources = make(map[string]bool)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func findInPaths(filename string, paths []string) string {
	for _, p := range paths {
		fullPath := filepath.Join(p, filename)
		if fileExists(fullPath) {
			return fullPath
		}
	}
	return ""
}

func findFileInRepo(filename string) string {
	var foundPath string
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == filename {
			foundPath = path
			return fmt.Errorf("Found") // Stop searching
		}
		return nil
	})
	if err != nil && err.Error() == "Found" {
		return foundPath
	}
	return ""
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
		targetDir = destInc
		checkPaths = includePaths
	} else {
		targetDir = destSrc
		checkPaths = sourcePaths
	}

	// Check if exists in paths
	existingPath := findInPaths(filename, checkPaths)
	if existingPath != "" {
		// Scan it for recursive includes/symbols
		scanFile(existingPath)
		return
	}
	
	// Not found locally, search in Repo
	path := findFileInRepo(filename)
	if path != "" {
		fmt.Printf("Found missing dependency: %s -> Copying to %s\n", filename, targetDir)
		copyFile(path, filepath.Join(targetDir, filename))
		
		// Scan the new file
		scanFile(filepath.Join(targetDir, filename))

		if isHeader {
			// Check for associated .c file
			cFile := strings.Replace(filename, ".h", ".c", 1)
			// Check if .c file already exists or processed
			if findInPaths(cFile, sourcePaths) == "" && !processedFiles[cFile] {
				// Search and copy associated .c file
				cPath := filepath.Join(filepath.Dir(path), cFile)
				if !fileExists(cPath) {
					cPath = findFileInRepo(cFile)
				}
				if cPath != "" {
					fmt.Printf("Found associated source: %s -> Copying to %s\n", cFile, destSrc)
					copyFile(cPath, filepath.Join(destSrc, cFile))
					processedFiles[cFile] = true
					scanFile(filepath.Join(destSrc, cFile))
				}
			}
		}
	}
}

func copyFile(src, dst string) {
	input, err := os.ReadFile(src)
	if err != nil {
		fmt.Println("Error reading", src)
		return
	}
	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		fmt.Println("Error writing", dst)
		return
	}
}

func RunDeps() {
	os.MkdirAll(destInc, 0755)
	os.MkdirAll(destSrc, 0755)
	// Seed with existing source files in all source paths
	for _, p := range sourcePaths {
		files, _ := filepath.Glob(filepath.Join(p, "*.c"))
		for _, f := range files {
			scanFile(f)
		}
	}
}
