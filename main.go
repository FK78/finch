package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

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
	if len(os.Args) < 2 {
		log.Fatal("Unknown command")
	}
	switch os.Args[1] {
	case "add":
		userHolding, err := getUserHolding()
		if err != nil {
			log.Fatal("Error:", err)
		}
		_, err = fetchSymbolDataAndCalculateProfitAndLossSinceBuy(userHolding, config)
		if err != nil {
			log.Fatal("Error:", err)
		}
		holdingsJson, err := loadHoldingsJSON()
		if err != nil {
			log.Fatal("Error:", err)
		}
		err = saveToHoldingsJSON(holdingsJson, userHolding)
		if err != nil {
			log.Fatal("Error:", err)
		}
	case "list":

	default:
		fmt.Println("unknown command")
	}
}
