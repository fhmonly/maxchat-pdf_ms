package main

import (
	"fmt"
	"log"
	"maxchat/pdf_ms/internal/config"
	"maxchat/pdf_ms/internal/router"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using system env")
	}

	db, err := config.NewMySQL()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer db.Close()

	envMode := os.Getenv("ENV")
	if envMode == "" {
		envMode = "development"
	}

	host := ":"
	if envMode == "development" || envMode == "d" || envMode == "dev" {
		host = "127.0.0.1"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := router.NewRouter(db)
	log.Printf("Server running on http://%s:%s\n", host, port)
	log.Fatal(http.ListenAndServe(host+":"+port, r))
}
