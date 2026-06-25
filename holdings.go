package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/tabwriter"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Holding struct {
	Ticker       string  `json:"Ticker"`
	BuyPrice     float64 `json:"BuyPrice"`
	AmountBought float64 `json:"AmountBought"`
}

func fetchSymbolDataAndCalculateProfitAndLossSinceBuy(userHolding Holding, finnhubToken Config) (float64, error) {
	req, err := http.NewRequest("GET", "https://finnhub.io/api/v1/quote?symbol="+userHolding.Ticker, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("X-Finnhub-Token", finnhubToken.Token)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	var data map[string]float64
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil {
		return 0, err
	}

	msg := message.NewPrinter(language.BritishEnglish)
	msg.Printf("Profit: %.2f\n", (data["c"]-userHolding.BuyPrice)*userHolding.AmountBought)
	return data["c"], nil
}

func getUserHolding() (Holding, error) {
	var userHolding Holding

	for i := 0; ; i++ {
		fmt.Print("Enter Ticker: ")
		_, err := fmt.Scan(
			&userHolding.Ticker,
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
			&userHolding.BuyPrice,
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
			&userHolding.AmountBought,
		)
		if err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}

	return userHolding, nil
}

func loadHoldingsJSON() ([]Holding, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return []Holding{}, err
	}
	holdingsFilePath := filepath.Join(userHomeDir, ".finch", "holdings.json")
	holdingsFile, err := os.ReadFile(holdingsFilePath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(holdingsFilePath, []byte("[]"), 0600)
		if err != nil {
			return []Holding{}, err
		}
		return []Holding{}, nil
	}
	var holdingsConfig []Holding
	if err := json.Unmarshal(holdingsFile, &holdingsConfig); err != nil {
		return []Holding{}, err
	}
	return holdingsConfig, nil
}

func saveToHoldingsJSON(holdings []Holding, userHolding Holding) error {
	holdingExists := false
	for i, holding := range holdings {
		if holding.Ticker == userHolding.Ticker {
			holdings[i].BuyPrice = (holding.AmountBought*holding.BuyPrice + userHolding.AmountBought*userHolding.BuyPrice) / (holding.AmountBought + userHolding.AmountBought)
			holdings[i].AmountBought += userHolding.AmountBought
			holdingExists = true
			break
		}
	}
	if !holdingExists {
		holdings = append(holdings, userHolding)
	}
	jsonData, err := json.Marshal(holdings)
	if err != nil {
		return err
	}
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	holdingsFilePath := filepath.Join(userHomeDir, ".finch", "holdings.json")

	err = os.WriteFile(holdingsFilePath, jsonData, 0600)
	if err != nil {
		return err
	}

	return nil
}

func displayHoldings(holdings []Holding) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "| Ticker\t| Price Per Share\t| Amount of Shares\t|")
	fmt.Fprintln(w, "| ------\t| ---------------\t| ----------------\t|")
	for _, holding := range holdings {
		fmt.Fprintf(w, "| %s\t| %.2f\t| %.2f\t|\n", holding.Ticker, holding.BuyPrice, holding.AmountBought)
	}
	w.Flush()
}
