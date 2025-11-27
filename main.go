package main

import (
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	maxChirpLength = 140
)

func main() {
	godotenv.Load()

	apiCfg := &apiConfig{
		fileServerHits: atomic.Int32{},
		JWT_SECRET:     os.Getenv("JWT_SECRET"),
		PLATFORM:       os.Getenv("PLATFORM"),
		API_KEY:        os.Getenv("POLKA_KEY"),
	}
	apiCfg.initDB()

	mux := http.NewServeMux()

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	apiCfg.registerRoutes(mux, handler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("could not start the server: %v", err)
	}

}
