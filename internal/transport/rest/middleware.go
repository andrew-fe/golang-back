package rest

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// statusRecorder envuelve http.ResponseWriter para capturar el status code
// efectivamente escrito, necesario para el middleware de logging.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// RequestLogger registra método, ruta, status y duración de cada petición.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", middleware.GetReqID(r.Context()),
		)
	})
}

// Recoverer captura panics en los handlers, los registra y responde 500
// en lugar de tumbar el servidor.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recuperado", "error", rec, "path", r.URL.Path)
				writeError(w, http.StatusInternalServerError, "error interno del servidor")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CORS habilita peticiones cross-origin básicas, configurable en el futuro
// vía variables de entorno si se necesita restringir orígenes.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// JSONContentType exige que las peticiones con body declaren
// Content-Type: application/json.
func JSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > 0 {
			ct := r.Header.Get("Content-Type")
			if ct != "application/json" && ct != "application/json; charset=utf-8" {
				writeError(w, http.StatusUnsupportedMediaType, "Content-Type debe ser application/json")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
