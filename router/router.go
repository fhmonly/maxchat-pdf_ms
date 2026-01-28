package router

import (
	"database/sql"
	"maxchat/pdf_ms/internal/module/pdf"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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
