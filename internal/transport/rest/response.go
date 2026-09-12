// Package rest expone la API HTTP de la aplicación (handlers, middlewares
// y el router), traduciendo entre el protocolo HTTP y la capa de servicio.
package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// errorBody es el envelope JSON usado para todas las respuestas de error.
type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// writeJSON serializa v como JSON con el status dado.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("rest: no se pudo codificar la respuesta JSON", "error", err)
	}
}

// writeError escribe un error JSON consistente con el status dado.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorBody{Error: errorDetail{Message: message}})
}

// writeValidationError escribe un 422 con el detalle de errores por campo.
func writeValidationError(w http.ResponseWriter, fields map[string]string) {
	writeJSON(w, http.StatusUnprocessableEntity, errorBody{
		Error: errorDetail{Message: "los datos proporcionados no son válidos", Fields: fields},
	})
}
