package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

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
