// Package model define las entidades de dominio de la aplicación.
package model

import (
	"strings"
	"time"
)

// Book representa un libro dentro de la biblioteca.
type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	ISBN      string    `json:"isbn,omitempty"`
	Year      int       `json:"year,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate aplica las reglas de negocio mínimas de un libro y devuelve
// un mapa de errores por campo cuando alguna no se cumple.
func (b *Book) Validate() map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(b.Title) == "" {
		errs["title"] = "el título es requerido"
	} else if len(b.Title) > 255 {
		errs["title"] = "el título no puede superar 255 caracteres"
	}

	if strings.TrimSpace(b.Author) == "" {
		errs["author"] = "el autor es requerido"
	} else if len(b.Author) > 255 {
		errs["author"] = "el autor no puede superar 255 caracteres"
	}

	if b.Year < 0 || b.Year > time.Now().Year()+1 {
		errs["year"] = "el año de publicación no es válido"
	}

	return errs
}

// Sanitize normaliza los campos de texto del libro (recorta espacios).
func (b *Book) Sanitize() {
	b.Title = strings.TrimSpace(b.Title)
	b.Author = strings.TrimSpace(b.Author)
	b.ISBN = strings.TrimSpace(b.ISBN)
}
