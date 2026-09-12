// Package store implementa la capa de persistencia sobre PostgreSQL.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/andrew-fe/golang-back/internal/model"
)

// BookStore expone las operaciones de persistencia disponibles para libros.
type BookStore struct {
	db *sql.DB
}

// NewBookStore crea un BookStore respaldado por la conexión dada.
func NewBookStore(db *sql.DB) *BookStore {
	return &BookStore{db: db}
}

// List devuelve una página de libros que coinciden con los parámetros de
// búsqueda, junto al total de coincidencias (sin paginar).
func (s *BookStore) List(ctx context.Context, params model.ListParams) ([]*model.Book, int, error) {
	total, err := s.count(ctx, params.Search)
	if err != nil {
		return nil, 0, err
	}

	q := `
		SELECT id, title, author, isbn, year, created_at, updated_at
		FROM books
		WHERE ($1 = '' OR title ILIKE '%' || $1 || '%' OR author ILIKE '%' || $1 || '%')
		ORDER BY id
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.QueryContext(ctx, q, params.Search, params.PageSize, params.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("store: listando libros: %w", err)
	}
	defer rows.Close()

	books := make([]*model.Book, 0, params.PageSize)
	for rows.Next() {
		b := &model.Book{}
		if err := scanBook(rows, b); err != nil {
			return nil, 0, fmt.Errorf("store: leyendo libro: %w", err)
		}
		books = append(books, b)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("store: iterando libros: %w", err)
	}

	return books, total, nil
}

func (s *BookStore) count(ctx context.Context, search string) (int, error) {
	q := `
		SELECT COUNT(*) FROM books
		WHERE ($1 = '' OR title ILIKE '%' || $1 || '%' OR author ILIKE '%' || $1 || '%')
	`

	var total int
	if err := s.db.QueryRowContext(ctx, q, search).Scan(&total); err != nil {
		return 0, fmt.Errorf("store: contando libros: %w", err)
	}
	return total, nil
}

// GetByID busca un libro por su identificador. Devuelve ErrNotFound si no existe.
func (s *BookStore) GetByID(ctx context.Context, id int) (*model.Book, error) {
	q := `
		SELECT id, title, author, isbn, year, created_at, updated_at
		FROM books WHERE id = $1
	`

	b := &model.Book{}
	err := scanBook(s.db.QueryRowContext(ctx, q, id), b)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("store: obteniendo libro %d: %w", id, err)
	}
	return b, nil
}

// Create inserta un nuevo libro y devuelve la entidad con su ID y timestamps.
func (s *BookStore) Create(ctx context.Context, book *model.Book) (*model.Book, error) {
	q := `
		INSERT INTO books (title, author, isbn, year)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, author, isbn, year, created_at, updated_at
	`

	created := &model.Book{}
	err := scanBook(
		s.db.QueryRowContext(ctx, q, book.Title, book.Author, book.ISBN, book.Year),
		created,
	)
	if err != nil {
		return nil, fmt.Errorf("store: creando libro: %w", err)
	}
	return created, nil
}

// Update actualiza un libro existente. Devuelve ErrNotFound si no existe.
func (s *BookStore) Update(ctx context.Context, id int, book *model.Book) (*model.Book, error) {
	q := `
		UPDATE books
		SET title = $1, author = $2, isbn = $3, year = $4, updated_at = now()
		WHERE id = $5
		RETURNING id, title, author, isbn, year, created_at, updated_at
	`

	updated := &model.Book{}
	err := scanBook(
		s.db.QueryRowContext(ctx, q, book.Title, book.Author, book.ISBN, book.Year, id),
		updated,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("store: actualizando libro %d: %w", id, err)
	}
	return updated, nil
}

// Delete elimina un libro por ID. Devuelve ErrNotFound si no existe.
func (s *BookStore) Delete(ctx context.Context, id int) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("store: eliminando libro %d: %w", id, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: verificando eliminación del libro %d: %w", id, err)
	}
	if affected == 0 {
		return ErrNotFound
	}

	return nil
}

// rowScanner abstrae *sql.Row y *sql.Rows para reutilizar el escaneo.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanBook(row rowScanner, b *model.Book) error {
	return row.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.Year, &b.CreatedAt, &b.UpdatedAt)
}
