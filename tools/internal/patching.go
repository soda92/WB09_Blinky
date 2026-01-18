package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func ApplyConfig(config ProjectConfig) {
	ApplyLPMSetting(config.DisableLPM)
	ApplyTraceSetting(config.EnableTrace)
	if config.MacAddress != "" {
		ApplyMacAddress(config.MacAddress)
	}
}

func ApplyLPMSetting(disable bool) {
	headerPath := filepath.Join("Core", "Inc", "app_conf.h")
	if !FileExists(headerPath) {
		return
	}

	content, err := os.ReadFile(headerPath)
	if err != nil {
		fmt.Printf("Error reading app_conf.h: %v\n", err)
		return
	}

	// Regex to find #define CFG_LPM_SUPPORTED (x)
	re := regexp.MustCompile(`(#define\s+CFG_LPM_SUPPORTED\s+)\(\d+\)`) // Note: The backslash before  in the original regex was unnecessary and has been removed.

	val := "1"
	if disable {
		val = "0"
	}

	newContent := re.ReplaceAllString(string(content), fmt.Sprintf("${1}(%s)", val))

	if string(content) != newContent {
		if err := os.WriteFile(headerPath, []byte(newContent), 0644); err != nil {
			fmt.Printf("Error writing app_conf.h: %v\n", err)
		} else {
			fmt.Printf("Applied Config: CFG_LPM_SUPPORTED set to (%s)\n", val)
		}
	}
}

func ApplyTraceSetting(enable bool) {
	headerPath := filepath.Join("Core", "Inc", "app_conf.h")
	if !FileExists(headerPath) {
		return
	}

	content, err := os.ReadFile(headerPath)
	if err != nil {
		fmt.Printf("Error reading app_conf.h: %v\n", err)
		return
	}

	re := regexp.MustCompile(`(#define\s+CFG_DEBUG_APP_TRACE\s+)\(\d+\)`)

	val := "0"
	if enable {
		val = "1"
	}

	newContent := re.ReplaceAllString(string(content), fmt.Sprintf("${1}(%s)", val))

	if string(content) != newContent {
		if err := os.WriteFile(headerPath, []byte(newContent), 0644); err != nil {
			fmt.Printf("Error writing app_conf.h: %v\n", err)
		} else {
			fmt.Printf("Applied Config: CFG_DEBUG_APP_TRACE set to (%s)\n", val)
		}
	}
}

func ApplyMacAddress(mac string) {
	headerPath := filepath.Join("Core", "Inc", "app_conf.h")
	if !FileExists(headerPath) {
		return
	}

	content, err := os.ReadFile(headerPath)
	if err != nil {
		fmt.Printf("Error reading app_conf.h: %v\n", err)
		return
	}

	// 1. Set MAC
	reMac := regexp.MustCompile(`(#define\s+CFG_PUBLIC_BD_ADDRESS\s+)\(0x[0-9a-fA-F]+\)`)
	newContent := reMac.ReplaceAllString(string(content), fmt.Sprintf("${1}(%s)", mac))

	// 2. Set Type to PUBLIC
	reType := regexp.MustCompile(`(#define\s+CFG_BD_ADDRESS_TYPE\s+)\w+`)
	newContent = reType.ReplaceAllString(newContent, "${1}HCI_ADDR_PUBLIC")

	if string(content) != newContent {
		if err := os.WriteFile(headerPath, []byte(newContent), 0644); err != nil {
			fmt.Printf("Error writing app_conf.h: %v\n", err)
		} else {
			fmt.Printf("Applied Config: MAC Address set to %s (Public)\n", mac)
		}
	}
}
