package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite"
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

func main() {
	http.HandleFunc("/cotacao", cotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func cotacaoHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log.Println("Request iniciado")
	defer log.Println("Request finalizado")
	select {
	case <-time.After(200 * time.Millisecond):
		res, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		var cotacao Cotacao
		err = json.Unmarshal(body, &cotacao)
		if err != nil {
			panic(err)
		}

		saveCotacao(cotacao)

		fmt.Println(cotacao.Usdbrl.Bid)

		// fmt.Println(string(body))
	case <-ctx.Done():
		http.Error(res, "Request canceelada", http.StatusRequestTimeout)
	}

}

func saveCotacao(cotacao Cotacao) {
	db, err := sql.Open("sqlite3", "./cotacao.db")
	if err != nil {
		panic(err)
	}

	result, err := db.Exec("INSERT INTO cotacao (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, createDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", cotacao.Usdbrl.Code, cotacao.Usdbrl.Codein, cotacao.Usdbrl.Name, cotacao.Usdbrl.High, cotacao.Usdbrl.Low, cotacao.Usdbrl.VarBid, cotacao.Usdbrl.PctChange, cotacao.Usdbrl.Bid, cotacao.Usdbrl.Ask, cotacao.Usdbrl.Timestamp, cotacao.Usdbrl.CreateDate)
	if err != nil {
		panic(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	fmt.Println("inserido", rowsAffected, "linha(s)")

}
