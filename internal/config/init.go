package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func SaveConfig(cfg Config) error {
	cfg.sanitize()

	path, err := DefaultConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return err
	}

	return nil
}
