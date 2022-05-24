package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
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
			break
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Println("=== " + url + " ===")
		fmt.Println(decoded)
	}
}

func sendJson(urlStr string, jsonBytes []byte, proxy string) {
	reader := bytes.NewReader(jsonBytes)
	req, err := http.NewRequest(
		"POST",
		urlStr,
		reader,
	)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	conf := &tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{
		TLSClientConfig: conf,
	}

	if len(proxy) > 0 {
		proxyUrl, _ := url.Parse(proxy)
		tr.Proxy = http.ProxyURL(proxyUrl)
	}

	client := &http.Client{
		Transport: tr,
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== " + urlStr + " ===")
	for key, val := range res.Header {
		fmt.Println(key + ":" + strings.Join(val, ","))
	}
	fmt.Println("")
	fmt.Println(beautifyJson(body))
}

func sendJsonFileViaProxy(url string, proxy string, jsonFile string) {
	bytes, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatal(err)
		return
	}

	sendJson(url, bytes, proxy)
}

func beautifyJson(jsonBytes []byte) string {
	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
	var decoded interface{}

	for {
		if err := decoder.Decode(&decoded); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		encoded, _ := json.MarshalIndent(decoded, "", "  ")
		return string(encoded)
	}

	return ""
}

func main() {
	if _, envErr := os.Stat(".env"); envErr == nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}

	//getJson(os.Getenv("JSON_URL"))

	params := map[string]interface{}{
		"string": "hello",
		"number": 123,
	}

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}

	sendJson(os.Getenv("ECHO_URL"), jsonBytes, "")

	if len(os.Args) == 2 {
		jsonFile := os.Args[1]
		sendJsonFileViaProxy(os.Getenv("API_URL"), os.Getenv("PROXY_URL"), jsonFile)
	}
}
