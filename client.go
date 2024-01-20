package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	requestTimeout = 300 * time.Millisecond
	requestUrl     = "http://localhost:8080/cotacao"
	fileName       = "cotacao.txt"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	r, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// Print to console
	log.Printf("Bid: %v\n", string(r))

	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	f.Write([]byte(fmt.Sprintf("Dólar: %v", string(r))))
}
