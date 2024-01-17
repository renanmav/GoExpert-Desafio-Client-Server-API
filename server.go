package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

const QuoteApiUrl = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

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

	//select {
	//case <-time.After(100 * time.Millisecond):
	//	fmt.Println("꩜ Request processed successfully")
	//	w.Write([]byte("🚀 Response successfully processed"))
	//case <-ctx.Done():
	//	fmt.Println("꩜ Request canceled by the client")
	//}
}

// TODO: add timeout
func FetchQuote(_ context.Context) (*Quote, error) {
	fmt.Println("꩜ Fetching quote...")

	req, err := http.Get(QuoteApiUrl)
	if err != nil {
		fmt.Println("꩜ Failed to fetch quote")
		return nil, err
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("꩜ Failed to read response")
		return nil, err
	}

	var data RawQuote
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Println("꩜ Failed to parse response")
		return nil, err
	}

	fmt.Println("꩜ Quote fetched successfully")
	return &data.Quote, nil
}
