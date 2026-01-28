package main

import (
	"fmt"
	"log"
	"maxchat/pdf_ms/internal/config"
	"maxchat/pdf_ms/router"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := router.NewRouter(db)
	log.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
