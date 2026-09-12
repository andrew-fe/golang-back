package store

import "errors"

// Errores de dominio devueltos por la capa de persistencia. Las capas
// superiores deben comparar con errors.Is en lugar de inspeccionar errores
// específicos del driver de base de datos.
var (
	// ErrNotFound indica que el recurso solicitado no existe.
	ErrNotFound = errors.New("recurso no encontrado")
)
