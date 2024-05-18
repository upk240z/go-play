package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type SimpleData struct {
	Code string `yaml:"code"`
	Age  int    `yaml:"age"`
}

func main() {
	filePath := "data/simple.yaml"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(file)
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var result SimpleData
	if err := decoder.Decode(&result); err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)

	structure := map[string]any{}
	yamlBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(yamlBytes, structure); err != nil {
		log.Fatal(err)
	}

	fmt.Println(structure)
}
