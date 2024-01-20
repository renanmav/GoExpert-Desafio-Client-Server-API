package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
}

type Quote struct {
	gorm.Model
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
	apiFetchTimeout = 200 * time.Millisecond
	dbWriteTimeout  = 10 * time.Millisecond
	apiUrl          = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	dbDsn           = "root:root@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	serverPath      = "/cotacao"
	serverPort      = ":8080"
)

func NewServer() *Server {
	db, err := gorm.Open(mysql.Open(dbDsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Quote{})
	return &Server{db}
}

func main() {
	log.Println("꩜ Initiating quote server...") // fmt or log?
	s := NewServer()

	http.HandleFunc(serverPath, s.QuoteHandler)
	http.ListenAndServe(serverPort, nil)
}

func (s *Server) QuoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("꩜ Request started")
	defer log.Println("꩜ Request finished")

	quote, err := s.FetchQuote(ctx)
	if err != nil {
		panic(err)
	}

	err = s.PersistQuote(ctx, quote)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(quote.Bid))
}

func (s *Server) FetchQuote(ctx context.Context) (*Quote, error) {
	log.Println("꩜ FetchQuote - Fetching quote...")
	startTime := time.Now()

	ctx, cancel := context.WithTimeout(ctx, apiFetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Println("꩜ FetchQuote - Failed to create GET request")
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("꩜ FetchQuote - Time elapsed: %v\n", time.Since(startTime))
		log.Println("꩜ FetchQuote - Failed to fetch quote")
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("꩜ FetchQuote - Failed to read response")
		return nil, err
	}

	var data RawQuote
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		log.Println("꩜ FetchQuote - Failed to parse response")
		return nil, err
	}

	log.Printf("꩜ FetchQuote - Time elapsed: %v\n", time.Since(startTime))
	log.Println("꩜ FetchQuote - Quote fetched successfully")
	return &data.Quote, nil
}

func (s *Server) PersistQuote(ctx context.Context, quote *Quote) error {
	log.Println("꩜ PersistQuote - Persisting quote...")
	startTime := time.Now()

	ctx, cancel := context.WithTimeout(ctx, dbWriteTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Create(&quote).Error
	if err != nil {
		log.Printf("꩜ PersistQuote - Time elapsed: %v\n", time.Since(startTime))
		log.Println("꩜ PersistQuote - Failed to persist quote")
		return err
	}

	log.Printf("꩜ PersistQuote - Time elapsed: %v\n", time.Since(startTime))
	log.Println("꩜ PersistQuote - Quote persisted successfully")
	return nil
}
