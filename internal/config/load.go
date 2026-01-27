package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func DefaultConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "printlock", "config.json"), nil
}

func LoadConfig() (*Config, error) {
	configFilePath, err := DefaultConfigPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		if cerr := file.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}
	defer cleanup()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	config.sanitize()
	return &config, nil
}
