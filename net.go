package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	var err error
	var res *http.Response
	var body []byte

	if _, envErr := os.Stat(".env"); envErr == nil {
		godotenv.Load(".env")
	}

	if res, err = http.Get(os.Getenv("JSON_URL")); err != nil {
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
			fmt.Println("EOF")
			break
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Println(decoded)
	}
}
