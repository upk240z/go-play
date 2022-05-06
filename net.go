package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func getJson(url string) {
	var err error
	var res *http.Response
	var body []byte

	if res, err = http.Get(url); err != nil {
		log.Fatal(err)
	}

	if body, err = io.ReadAll(res.Body); err != nil {
		log.Fatal(err)
	}

	jsonStr := string(body)
	decoder := json.NewDecoder(strings.NewReader(jsonStr))
	var decoded interface{}

	for {
		if err = decoder.Decode(&decoded); err == io.EOF {
			fmt.Println("@@@@@ JSON fin @@@@")
			break
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Println(decoded)
	}
}

func sendJson(url string, params interface{}) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(jsonData)
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		url,
		reader,
	)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}

func main() {
	if _, envErr := os.Stat(".env"); envErr == nil {
		godotenv.Load(".env")
	}

	getJson(os.Getenv("JSON_URL"))
	sendJson(os.Getenv("ECHO_URL"), map[string]interface{}{
		"string": "hello",
		"number": 123,
	})
}
