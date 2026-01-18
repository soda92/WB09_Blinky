package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const repoPath = "/home/soda/STM32Cube/Repository/STM32Cube_FW_WB0_V1.4.0"
const destInc = "Core/Inc"
const destSrc = "Core/Src"

// Directories to check before deciding a file is missing
var includePaths = []string{
	"Core/Inc",
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
	"STM32_BLE/App",
	"STM32_BLE/Target",
	"Drivers/STM32WB0x_HAL_Driver/Src",
	"Drivers/BSP/STM32WB0x-nucleo",
	"Middlewares/ST/STM32_BLE/evt_handler/src",
	"Middlewares/ST/STM32_BLE/stack/config",
}

var processedHeaders = make(map[string]bool)
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
		if strings.HasPrefix(line, "#include") {
			parts := strings.Split(line, "\"")
			if len(parts) > 1 {
				header := parts[1]
				processHeader(header)
			}
		}
	}
}

func processHeader(header string) {
	if processedHeaders[header] {
		return
	}
	processedHeaders[header] = true // Mark processed early to avoid cycles

	// Check if exists in include paths
	existingPath := findInPaths(header, includePaths)
	if existingPath != "" {
		// Scan it for recursive includes
		scanFile(existingPath)
		return
	}
	
	// Not found locally, search in Repo
	path := findFileInRepo(header)
	if path != "" {
		fmt.Printf("Found missing header: %s -> Copying to %s\n", header, destInc)
		copyFile(path, filepath.Join(destInc, header))
		
		// Scan the new header
		scanFile(filepath.Join(destInc, header))

		// Check for associated .c file
		cFile := strings.Replace(header, ".h", ".c", 1)
		
		// Check if .c file already exists in project source paths
		if findInPaths(cFile, sourcePaths) != "" {
			return
		}

		// Not found locally, search in Repo
		// Check same dir as header first
		cPath := filepath.Join(filepath.Dir(path), cFile)
		if !fileExists(cPath) {
			// Try finding globally in Repo
			cPath = findFileInRepo(cFile)
		}
		
		if cPath != "" && !foundSources[cFile] {
			fmt.Printf("Found associated source: %s -> Copying to %s\n", cFile, destSrc)
			copyFile(cPath, filepath.Join(destSrc, cFile))
			foundSources[cFile] = true
			// Scan the new source
			scanFile(filepath.Join(destSrc, cFile))
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
	// Seed with existing source files in all source paths
	for _, p := range sourcePaths {
		files, _ := filepath.Glob(filepath.Join(p, "*.c"))
		for _, f := range files {
			scanFile(f)
		}
	}
}