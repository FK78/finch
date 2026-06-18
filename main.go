package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Holding struct {
	ticker       string
	buyPrice     float64
	amountBought float64
}

func main() {
	finnhubToken := config().Token

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
