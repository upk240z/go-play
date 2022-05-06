package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

func main() {
	file, _ := os.Open("data/sample.yaml")
	decoder := yaml.NewDecoder(file)
	var result map[string]interface{}
	decoder.Decode(&result)
	file.Close()
	fmt.Println(result)
}
