package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/adi290491/chirpy/internal/database"
)

var (
	dbURL string
)

func InitDB(c *apiConfig) {
	dbURL = os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	c.db = dbQueries

}
