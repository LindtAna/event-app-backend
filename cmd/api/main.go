package main

import (
	"database/sql"
	"log"
)

func main() {
	db, err := sql.Open("sqlite", "./data.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}
