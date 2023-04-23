package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"play-ground/core"
	"time"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal(err)
		}
	}

	db := core.NewDatabase("mysql", os.Getenv("MYSQL_DSN"))

	rows := db.All(
		`SELECT * FROM zipcode WHERE zipcode LIKE :zip`,
		map[string]interface{}{
			"zip": "498%",
		},
	)

	for _, columns := range rows {
		fmt.Println(*columns["prefecture_id"], *columns["city"], *columns["town"])
	}

	row := db.Row(
		`SELECT * FROM zipcode WHERE zipcode = :zip`,
		map[string]interface{}{
			"zip": "0640918",
		},
	)

	for key, val := range row {
		fmt.Println(key, *val)
	}

	parameters := map[string]interface{}{
		"value": time.Now().Format("2006-01-02 15:04:05"),
	}
	affected := db.Exec(`INSERT INTO foo (bar) VALUES (:value)`, parameters)
	fmt.Printf("affected: %d\n", affected)

	db.Close()
}
