package service

import "errors"

// ErrNotFound indica que el recurso solicitado no existe.
var ErrNotFound = errors.New("recurso no encontrado")

// ValidationError agrupa errores de validación por campo. La capa de
// transporte puede inspeccionar Fields para construir una respuesta detallada.
type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return "los datos proporcionados no son válidos"
}

// NewValidationError construye un ValidationError a partir de un mapa de
// errores por campo. Devuelve nil si el mapa está vacío.
func NewValidationError(fields map[string]string) error {
	if len(fields) == 0 {
		return nil
	}
	return &ValidationError{Fields: fields}
}
