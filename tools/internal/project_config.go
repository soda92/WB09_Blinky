package internal

import (
	"encoding/json"
	"os"
)

type ProjectConfig struct {
	DisableLPM bool `json:"disable_lpm"`
}

const ConfigFile = "wb09.json"

func LoadConfig() (ProjectConfig, error) {
	var config ProjectConfig
	if !FileExists(ConfigFile) {
		return config, nil // Default zero values (false)
	}
	bytes, err := os.ReadFile(ConfigFile)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(bytes, &config)
	return config, err
}

func SaveConfig(config ProjectConfig) error {
	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile, bytes, 0644)
}
