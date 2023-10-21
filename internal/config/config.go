package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	RootMap map[string]string `json:"root_map"`
}

func ParseConfig(configPath string) (*Config, error) {
	byteValue, err := os.ReadFile(configPath)

	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{}

	if err := json.Unmarshal(byteValue, config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config file: %w", err)
	}

	return config, nil
}
