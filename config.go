package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Token string `json:"finnhubToken"`
}

func ensureFinchDirExists() error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	finchDir := filepath.Join(userHomeDir, ".finch")

	return os.MkdirAll(finchDir, 0755)
}

func loadConfig() (Config, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	finchConfigFile, err := os.ReadFile(filepath.Join(userHomeDir, ".finch", "config.json"))
	if err != nil {
		return Config{}, err
	}
	var finchConfig Config
	if err := json.Unmarshal(finchConfigFile, &finchConfig); err != nil {
		return Config{}, err
	}

	return finchConfig, nil
}

func saveConfig(config Config) error {
	jsonData, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(userHomeDir, ".finch", "config.json")

	err = os.WriteFile(configFilePath, jsonData, 0600)
	if err != nil {
		return err
	}

	return nil
}

func promptForAPIKey() (Config, error) {
	config := Config{}
	fmt.Print("Enter your Finnhub API Key: ")
	_, err := fmt.Scan(&config.Token)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
