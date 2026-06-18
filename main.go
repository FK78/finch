package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Config struct {
	Token string `json:"finnhubToken"`
}
type Holding struct {
	ticker       string
	buyPrice     float64
	amountBought float64
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

func addHolding() Holding {
	var userHolding Holding

	for i := 0; ; i++ {
		fmt.Print("Enter ticker: ")
		_, err := fmt.Scan(
			&userHolding.ticker,
		)
		if err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}

	for i := 0; ; i++ {
		fmt.Print("Enter buy price: ")
		_, err := fmt.Scan(
			&userHolding.buyPrice,
		)
		if err != nil {
			fmt.Println(err)
		} else {
			break
		}

	}

	for i := 0; ; i++ {
		fmt.Print("Enter amount bought: ")
		_, err := fmt.Scan(
			&userHolding.amountBought,
		)
		if err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}

	return userHolding
}

func main() {
	if err := ensureFinchDirExists(); err != nil {
		log.Fatal(err)
	}
	config, err := loadConfig()
	if errors.Is(err, os.ErrNotExist) {
		config, err = promptForAPIKey()
		if err != nil {
			log.Fatal(err)
		}
		saveConfig(config)
	}
	finnhubToken := config.Token

	userHolding := addHolding()
	req, err := http.NewRequest("GET", "https://finnhub.io/api/v1/quote?symbol="+userHolding.ticker, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("X-Finnhub-Token", finnhubToken)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	var data map[string]float64
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil {
		fmt.Println(err)
		return
	}

	msg := message.NewPrinter(language.BritishEnglish)
	msg.Printf("Profit: %.2f\n", (data["c"]-userHolding.buyPrice)*userHolding.amountBought)
}
