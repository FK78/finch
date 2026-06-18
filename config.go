package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type FinnHubToken struct {
	Token string `json:"finnhubToken"`
}

func config() FinnHubToken {
	var tokenConfig FinnHubToken
	userHomeDir, _ := os.UserHomeDir()
	finchConfigDir := userHomeDir + "/.finch/"
	if _, err := os.Stat(finchConfigDir); err == nil {
		file, err := os.ReadFile(finchConfigDir + "/.config")
		check(err)

		if err := json.Unmarshal(file, &tokenConfig); err != nil {
			log.Fatalf("Error unmarshalling JSON: %v", err)
		}

		return tokenConfig
	} else if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(finchConfigDir, 0755)
		check(err)
		file, err := os.Create(finchConfigDir + "/.config")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return FinnHubToken{}
		}
		defer file.Close()

		fmt.Println("Enter your Finnhub Token: ")
		_, err = fmt.Scan(&tokenConfig.Token)
		check(err)

		jsonData, err := json.Marshal(tokenConfig)
		check(err)

		_, err = file.Write(jsonData)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return tokenConfig
		}
		fmt.Println("Finnhub Token Saved")
		return tokenConfig
	} else {
		return FinnHubToken{}
	}
}
