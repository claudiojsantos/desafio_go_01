package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	checkError(err)

	res, err := http.DefaultClient.Do(req)
	checkError(err)
	defer res.Body.Close()

	var data map[string]string
	err = json.NewDecoder(res.Body).Decode(&data)
	checkError(err)

	bid, ok := data["bid"]
	if !ok {
		panic("bid not found")
	}

	saveFile(bid)
}

func saveFile(valor string) {
	file, err := os.Create("cotacao.txt")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.Write([]byte("Dolar: {" + valor + "}"))
	if err != nil {
		panic(err)
	}

	return
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
