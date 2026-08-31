package main

import (
	"database/sql"
	"event-app/internal/database"
	"event-app/internal/env"
	"log"

	_ "github.com/joho/godotenv/autoload"
	_ "modernc.org/sqlite"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	db, err := sql.Open("sqlite", "./data.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "somes-secret-1234"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
