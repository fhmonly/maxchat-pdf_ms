package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fhmonly/maxchat-pdfms/internal/config"
	"github.com/fhmonly/maxchat-pdfms/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	isProd := func() bool {
		switch env {
		case "production", "prod", "p":
			return true
		default:
			return false
		}
	}()

	if !isProd {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	db, err := config.NewMySQL()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer db.Close()

	host := ":"
	if !isProd {
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
