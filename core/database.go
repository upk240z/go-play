package core

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"regexp"
	"strings"
	"time"
)

type Database struct {
	driver  string
	dsn     string
	db      *sql.DB
	result  sql.Result
	query   string
	named   map[string]any
	unnamed []any
}

func NewDatabase(driver, dsn string) *Database {
	instance := Database{driver: driver, dsn: dsn}
	return &instance
}

func (database *Database) connect() {
	if database.db != nil {
		return
	}

	if db, err := sql.Open(database.driver, database.dsn); err != nil {
		log.Fatal(err)
	} else {
		database.db = db
	}

	database.db.SetMaxOpenConns(0)
	database.db.SetMaxIdleConns(10)
	database.db.SetConnMaxLifetime(time.Minute * 5)
}

func (database *Database) anonymize() {
	pattern, _ := regexp.Compile(`:([a-z\d\-_]+)`)

	database.unnamed = []any{}

	for {
		if !pattern.MatchString(database.query) {
			break
		}

		mark := pattern.FindString(database.query)
		p, exists := database.named[mark[1:]]
		if !exists {
			log.Fatal("parameter not found: " + mark[1:])
		}

		database.unnamed = append(database.unnamed, p)
		database.query = strings.Replace(database.query, mark, "?", 1)
	}
}

func (database *Database) doCommon(query string, parameters map[string]any) {
	database.connect()
	database.query = query
	database.named = parameters
	database.anonymize()
}

func (database *Database) All(query string, parameters map[string]any) []map[string]*string {
	database.doCommon(query, parameters)

	rows, err := database.db.Query(database.query, database.unnamed...)
	if err != nil {
		log.Fatal(err)
	}

	columns, _ := rows.Columns()

	var results []map[string]*string

	for rows.Next() {
		var pointers []any
		mapValues := map[string]*string{}
		for _, columnName := range columns {
			var col string
			pointers = append(pointers, &col)
			mapValues[columnName] = &col
		}

		if err := rows.Scan(pointers...); err != nil {
			log.Println(err)
			return results
		}

		results = append(results, mapValues)
	}

	return results
}

func (database *Database) Row(query string, parameters map[string]any) map[string]*string {
	rows := database.All(query, parameters)

	if len(rows) == 0 {
		return map[string]*string{}
	}

	return rows[0]
}

func (database *Database) Exec(query string, parameters map[string]any) int64 {
	database.doCommon(query, parameters)

	if result, err := database.db.Exec(database.query, database.unnamed...); err != nil {
		log.Fatal(err)
	} else {
		database.result = result
	}

	affected, err := database.result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	return affected
}

func (database *Database) LastInsertId() int64 {
	if database.result == nil {
		return 0
	}

	lastId, err := database.result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return lastId
}

func (database *Database) Close() {
	if database.db == nil {
		return
	}

	if err := database.db.Close(); err != nil {
		log.Fatal(err)
	}
}

func (database *Database) Begin() *sql.Tx {
	tx, err := database.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	return tx
}

func (database *Database) Prepare(query string) *sql.Stmt {
	stmt, err := database.db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}
