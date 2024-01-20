package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
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
	apiUrl          = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	apiFetchTimeout = 200 * time.Millisecond
	dbWriteTimeout  = 10 * time.Millisecond
)

func main() {
	log.Println("꩜ Initiating quote server...") // fmt or log?
	http.HandleFunc("/cotacao", QuoteHandler)
	http.ListenAndServe(":8080", nil)
}

func QuoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("꩜ Request started")
	defer log.Println("꩜ Request finished")

	quote, err := FetchQuote(ctx)
	if err != nil {
		panic(err)
	}

	// TODO: persist on database with timeout of 10ms

	result, err := json.Marshal(quote)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func FetchQuote(ctx context.Context) (*Quote, error) {
	log.Println("꩜ FetchQuote - Fetching quote...")
	startTime := time.Now()

	// Timeout with ctx vs http client?
	ctx, cancel := context.WithTimeout(ctx, apiFetchTimeout)
	defer cancel()
	client := http.Client{
		Timeout: apiFetchTimeout,
	}

	req, err := client.Get(apiUrl)
	if err != nil {
		log.Println("꩜ FetchQuote - Failed to fetch quote")
		return nil, err
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("꩜ FetchQuote - Failed to read response")
		return nil, err
	}

	var data RawQuote
	err = json.Unmarshal(res, &data)
	if err != nil {
		log.Println("꩜ FetchQuote - Failed to parse response")
		return nil, err
	}

	// TODO: log time elapsed between start and end of request
	log.Printf("꩜ FetchQuote - Time elapsed: %v\n", time.Since(startTime))
	log.Println("꩜ FetchQuote - Quote fetched successfully")
	return &data.Quote, nil
}
