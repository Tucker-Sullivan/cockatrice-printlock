package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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

func (c *Config) sanitize() {
	c.CardsXMLPath = strings.TrimSpace(c.CardsXMLPath)
	c.DeckDir = strings.TrimSpace(c.DeckDir)
	cleanedSets := make([]string, 0, len(c.GlobalSetPriority))
	for _, setCode := range c.GlobalSetPriority {
		tSetCode := strings.TrimSpace(setCode)
		if tSetCode == "" {
			continue
		}
		cleanedSets = append(cleanedSets, tSetCode)
	}
	c.GlobalSetPriority = cleanedSets
}
