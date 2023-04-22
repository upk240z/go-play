package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"play-ground/core"
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
		`SELECT * FROM address WHERE zipcode = :zip`,
		map[string]interface{}{
			"zip": "0640918",
		},
	)

	for key, val := range row {
		fmt.Println(key, *val)
	}

	db.Close()
}
