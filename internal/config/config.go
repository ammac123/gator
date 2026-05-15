package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	configFileName = ".gatorconfig.json"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home + "/" + configFileName, nil
}

func write(cfg Config) error {
	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	fp, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(fp, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Read() Config {
	path, err := getConfigFilePath()
	if err != nil {
		s := fmt.Errorf("Could not locate home dir: %w", err)
		fmt.Print(s)
		return Config{}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		s := fmt.Errorf("Error reading data: %w", err)
		fmt.Print(s)
		return Config{}
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		s := fmt.Errorf("Could not parse JSON: %w", err)
		fmt.Print(s)
		return Config{}
	}

	return config
}

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUsername = user
	err := write(*cfg)
	if err != nil {
		return fmt.Errorf("Could not set user: %w", err)
	}
	return nil
}
