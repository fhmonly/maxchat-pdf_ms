package router

import (
	"database/sql"
	"net/http"

	"github.com/fhmonly/maxchat-pdfms/internal/module/pdf"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // seconds
	}))

	pdfRepo := pdf.NewRepository(db)
	pdfService := pdf.NewService(pdfRepo, "./uploads/pdf")
	pdfHandler := pdf.NewHandler(pdfService)

	r.Route("/api", func(r chi.Router) {
		r.Route("/pdf", func(r chi.Router) {
			r.Post("/generate", pdfHandler.Generate)
			r.Post("/upload", pdfHandler.Upload)
			r.Get("/list", pdfHandler.List)
			r.Delete("/{id}", pdfHandler.Delete)
		})
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return r
}
