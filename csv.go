package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, _ := os.Open("data/sample.csv")
	reader := csv.NewReader(file)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Println(line)
	}
	err := file.Close()
	if err != nil {
		log.Fatal(err)
	}
}
