package rest

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// NewRouter construye el http.Handler completo de la API, incluyendo
// middlewares globales, health checks y las rutas de negocio versionadas.
func NewRouter(bookService BookService, db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(Recoverer)
	r.Use(RequestLogger)
	r.Use(CORS)
	r.Use(chimw.Timeout(30 * time.Second))

	health := NewHealthHandler(db)
	r.Get("/healthz", health.Live)
	r.Get("/readyz", health.Ready)

	books := NewBookHandler(bookService)
	r.Route("/api/v1/books", func(r chi.Router) {
		r.Use(JSONContentType)
		r.Get("/", books.List)
		r.Post("/", books.Create)
		r.Get("/{id}", books.Get)
		r.Put("/{id}", books.Update)
		r.Delete("/{id}", books.Delete)
	})

	return r
}
