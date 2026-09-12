package rest

import (
	"database/sql"
	"net/http"
)

// HealthHandler expone endpoints de liveness y readiness para orquestadores
// (Kubernetes, Docker, balanceadores de carga, etc.).
type HealthHandler struct {
	db *sql.DB
}

// NewHealthHandler crea un HealthHandler que verifica la conexión db en Ready.
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Live responde 200 mientras el proceso esté vivo, sin depender de nada externo.
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready responde 200 solo si las dependencias externas (la base de datos)
// están disponibles.
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	if err := h.db.PingContext(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
