package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

		PnL, err := fetchSymbolDataAndCalculateProfitAndLossSinceBuy(userHolding, config)
		if err != nil {
			log.Fatal("Error:", err)
		}
		msg := message.NewPrinter(language.BritishEnglish)
		msg.Printf("Profit/Loss: %.2f\n", (PnL))

		holdingsJson, err := loadHoldingsJSON()
		if err != nil {
			log.Fatal("Error:", err)
		}

		err = saveToHoldingsJSON(holdingsJson, userHolding)
		if err != nil {
			log.Fatal("Error:", err)
		}
	case "list":
		holdingsJson, err := loadHoldingsJSON()
		if err != nil {
			log.Fatal("Error:", err)
		}
		displayHoldings(holdingsJson, config)
	default:
		fmt.Println("unknown command")
	}
}
