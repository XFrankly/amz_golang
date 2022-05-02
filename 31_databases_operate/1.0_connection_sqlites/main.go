package main

import (
	"database/sql"

	sqlite3 "github.com/mattn/go-sqlite3"
)

var (
	Logg := log.New(os.Stderr, "INFO -- ", 13)
)

func run() error {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}
	defer db.Close()
	// do database sturr herr
	Logg.Println(sqlite3) 
}

func main() {
	run()
}

