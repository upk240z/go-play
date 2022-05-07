package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
	"play-ground/core"
)

func main() {
	if _, envErr := os.Stat(".env"); envErr == nil {
		godotenv.Load(".env")
	}

	db := core.NewDatabase("mysql", os.Getenv("MYSQL_DSN"))

	rows := db.All(
		`SELECT * FROM zipcode WHERE zipcode LIKE ?`,
		"152%",
	)

	for _, columns := range rows {
		fmt.Println(*columns["prefecture_id"], *columns["city"], *columns["town"])
	}

	db.Close()
}
