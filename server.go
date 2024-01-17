package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Quote struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type RawQuote struct {
	Quote Quote `json:"USDBRL"`
}

const (
	QuoteApiUrl          = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	QuoteApiFetchTimeout = 700 * time.Millisecond // 200ms will fail
)

func main() {
	fmt.Println("꩜ Initiating quote server...") // fmt or log?
	http.HandleFunc("/cotacao", QuoteHandler)
	http.ListenAndServe(":8080", nil)
}

func QuoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Println("꩜ Request started")
	defer fmt.Println("꩜ Request finished")

	quote, err := FetchQuote(ctx)
	if err != nil {
		panic(err)
	}

	result, err := json.Marshal(quote)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func FetchQuote(ctx context.Context) (*Quote, error) {
	fmt.Println("꩜ FetchQuote - Fetching quote...")

	// Timeout with ctx vs http client?
	ctx, cancel := context.WithTimeout(ctx, QuoteApiFetchTimeout)
	defer cancel()
	client := http.Client{
		Timeout: QuoteApiFetchTimeout,
	}

	req, err := client.Get(QuoteApiUrl)
	if err != nil {
		fmt.Println("꩜ FetchQuote - Failed to fetch quote")
		return nil, err
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("꩜ FetchQuote - Failed to read response")
		return nil, err
	}

	var data RawQuote
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Println("꩜ FetchQuote - Failed to parse response")
		return nil, err
	}

	fmt.Println("꩜ FetchQuote - Quote fetched successfully")
	return &data.Quote, nil
}
