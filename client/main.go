package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
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

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	var cotacao Cotacao
	err = json.NewDecoder(res.Body).Decode(&cotacao)
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, res.Body)

	// saveFile(cotacao)
}

func saveFile(cotacao Cotacao) {
	file, err := os.Create("cotacao.txt")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.Write([]byte(cotacao.Code))
	if err != nil {
		panic(err)
	}
}
