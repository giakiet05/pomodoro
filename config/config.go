package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Work       int `json:"work"`        // in minutes
	ShortBreak int `json:"short_break"` // in minutes
	LongBreak  int `json:"long_break"`  // in minutes
}

func DefaultConfig() *Config {
	return &Config{
		Work:       25,
		ShortBreak: 5,
		LongBreak:  15,
	}
}

func SaveConfig(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // Return default config if file doesn't exist
		}
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
