package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ServerURL    string `json:"server_url"`
	DefaultUser  string `json:"default_user"`
	ItemsPerPage int    `json:"items_per_page"`
}

func Load() Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return defaultConfig()
	}

	configPath := filepath.Join(homeDir, ".config", "jellyfin-tui", "config.json")
	file, err := os.Open(configPath)
	if err != nil {
		return defaultConfig()
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return defaultConfig()
	}

	return config
}

func defaultConfig() Config {
	return Config{
		ServerURL:    "http://localhost:8096",
		DefaultUser:  "",
		ItemsPerPage: 20,
	}
}

func Save(config Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "jellyfin-tui")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}
