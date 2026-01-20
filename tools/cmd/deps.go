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
	"UTIL_PowerDriver":     "stm32_lpm_if.c",
	"HOST_TO_LE_16":        "ble_types.h", // Just in case
}

var requiredLibs = []string{
	"Middlewares/ST/STM32_BLE/stack/lib/stm32wb0x_ble_stack.a",
	"Middlewares/ST/STM32_BLE/cryptolib/libcrypto.a",
}

func runDeps() {
	os.MkdirAll(internal.LibIncDir, 0755)
	os.MkdirAll(internal.LibSrcDir, 0755)

	// 1. Copy Static Libraries (Binaries)
	fmt.Println("Checking static libraries...")
	for _, libRelPath := range requiredLibs {
		localPath := libRelPath // Same relative path locally
		if !internal.FileExists(localPath) {
			sdkPath := filepath.Join(internal.SDKPath, libRelPath)
			if internal.FileExists(sdkPath) {
				fmt.Printf("Copying library: %s\n", libRelPath)
				if err := internal.CopyFile(sdkPath, localPath); err != nil {
					fmt.Printf("Error copying library %s: %v\n", libRelPath, err)
				}
			} else {
				fmt.Printf("Warning: Library not found in SDK: %s\n", sdkPath)
			}
		}
	}

	// Seed with existing source files in all source paths
	for _, p := range sourcePaths {
		files, _ := filepath.Glob(filepath.Join(p, "*.c"))
		for _, f := range files {
			scanFile(f)
		}
	}

	// Patch osal.c if it exists
	osalPath := filepath.Join(internal.LibSrcDir, "osal.c")
	if internal.FileExists(osalPath) {
		content := "\nvoid Osal_MemCpy(void *dest, const void *src, unsigned int size) {\n    memcpy(dest, src, size);\n}\n\nvoid Osal_MemSet(void *ptr, int value, unsigned int size) {\n    memset(ptr, value, size);\n}\n\nint Osal_MemCmp(const void *ptr1, const void *ptr2, unsigned int size) {\n    return memcmp(ptr1, ptr2, size);\n}\n"
		// Check if already patched
		existing, _ := os.ReadFile(osalPath)
		if !strings.Contains(string(existing), "Osal_MemCpy(void") {
			internal.AppendToFile(osalPath, content)
			fmt.Println("Patched osal.c with missing functions.")
		}
	}

	// Patch cpu_context_switch.s
	cpuPath := filepath.Join(internal.LibSrcDir, "cpu_context_switch.s")
	if internal.FileExists(cpuPath) {
		internal.ReplaceInFile(cpuPath, "../Modules/asm.h", "asm.h")
		fmt.Println("Patched cpu_context_switch.s include path.")
	}

	// Always check/copy cryptolib_hw_aes.c as it is needed by libcrypto.a
	processDependency("cryptolib_hw_aes.c")
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
				// fmt.Printf("DEBUG: Found symbol %s in %s -> triggering %s\n", symbol, path, filename)
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

	// Filter DTM sources (keep headers)
	if !isHeader && (strings.Contains(filename, "dtm_") || strings.Contains(filename, "transport_layer") || strings.Contains(filename, "hci_parser")) {
		return
	}

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
			processDependency(cFile)
		}
	}
}
