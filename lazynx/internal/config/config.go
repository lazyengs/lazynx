package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	AppName    = ".lazynx"
	ConfigName = "lazynx.json"
)

type Config struct {
	Logs string `json:"logs"`
}

func new() *Config {
	return &Config{
		Logs: getDefaultLogFile(),
	}
}

func LoadConfiguration() *Config {
	config := new()

	homeDir, err := getHomeDir()
	if err != nil {
		return config
	}

	globalConfigPath := filepath.Join(homeDir, AppName, ConfigName)
	if globalConfig, eC := loadConfigFromFile(globalConfigPath); eC == nil {
		config = config.overrideWith(globalConfig)
	}

	return config
}

func loadConfigFromFile(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) overrideWith(target *Config) *Config {
	result := *c

	if target != nil {
		if target.Logs != "" {
			result.Logs = target.Logs
		}
	}

	return &result
}

func getHomeDir() (string, error) {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return dir, nil
	}

	if dir, err := os.UserHomeDir(); err == nil {
		return dir, nil
	}

	return "", errors.New("unable to determine user config directory")
}

func getDefaultLogFile() string {
	homeDir, _ := getHomeDir()
	return filepath.Join(homeDir, AppName, "logs", "lazynx.log")
}
