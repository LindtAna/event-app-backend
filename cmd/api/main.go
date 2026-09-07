package main

import (
	"database/sql"
	_ "event-app/docs"
	"event-app/internal/database"
	"event-app/internal/env"
	"log"

	_ "github.com/joho/godotenv/autoload"
	_ "modernc.org/sqlite"
)

// @title Go gin REST API
// @version 1.0
// @description a REST API in Go using Gin framework
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer tken in the forman **Bearer &lt;token&gt;**
type application struct {
	port         int
	jwtSecret    string
	clientOrigin string
	models       database.Models
}

func main() {
	db, err := sql.Open("sqlite", "./data.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port:         env.GetEnvInt("PORT", 8080),
		jwtSecret:    env.GetEnvString("JWT_SECRET", "somes-secret-1234"),
		clientOrigin: env.GetEnvString("CLIENT_ORIGIN", "http://localhost:5173"),
		models:       models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
