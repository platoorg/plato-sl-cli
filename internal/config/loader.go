package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads and parses a platosl.yaml configuration file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s\n\nRun 'platosl init' to create a new configuration", path)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults
	if cfg.Version == "" {
		cfg.Version = "v1"
	}
	if len(cfg.Schemas) == 0 {
		cfg.Schemas = []string{"schemas/"}
	}
	if cfg.Generate == nil {
		cfg.Generate = make(map[string]GenConfig)
	}

	return &cfg, nil
}

// Save writes a configuration to a file
func Save(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Exists checks if a config file exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
