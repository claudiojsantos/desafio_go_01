package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	Usdbrl struct {
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
	} `json:"USDBRL"`
}

type CotacaoBid struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", cotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func cotacaoHandler(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	select {
	case <-time.After(200 * time.Millisecond):
		res, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
		checkErr(err)
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		checkErr(err)

		var cotacao Cotacao
		err = json.Unmarshal(body, &cotacao)
		checkErr(err)

		saveCotacao(cotacao)

		resp.Header().Set("Content-Type", "application/json")

		BidCotacao := CotacaoBid{Bid: cotacao.Usdbrl.Bid}
		err = json.NewEncoder(resp).Encode(BidCotacao)
		checkErr(err)
	case <-ctx.Done():
		http.Error(resp, "Request canceelada", http.StatusRequestTimeout)
	}
}

func saveCotacao(cotacao Cotacao) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	db, err := sql.Open("sqlite3", "./cotacao.db")
	checkErr(err)

	stmt, err := db.PrepareContext(ctx, "INSERT INTO cotacao (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, createDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	checkErr(err)

	stmt.Exec(cotacao.Usdbrl.Code, cotacao.Usdbrl.Codein, cotacao.Usdbrl.Name, cotacao.Usdbrl.High, cotacao.Usdbrl.Low, cotacao.Usdbrl.VarBid, cotacao.Usdbrl.PctChange, cotacao.Usdbrl.Bid, cotacao.Usdbrl.Ask, cotacao.Usdbrl.Timestamp, cotacao.Usdbrl.CreateDate)

	defer db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
		return
	}
}
