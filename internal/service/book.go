// Package service contiene la lógica de negocio de la aplicación,
// independiente de los detalles de transporte (HTTP) y persistencia (SQL).
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/andrew-fe/golang-back/internal/model"
	"github.com/andrew-fe/golang-back/internal/store"
)

// BookStore es el contrato de persistencia que necesita BookService.
// Se define en el consumidor (no en la capa de store) siguiendo el
// principio de que las interfaces se declaran donde se usan.
type BookStore interface {
	List(ctx context.Context, params model.ListParams) ([]*model.Book, int, error)
	GetByID(ctx context.Context, id int) (*model.Book, error)
	Create(ctx context.Context, book *model.Book) (*model.Book, error)
	Update(ctx context.Context, id int, book *model.Book) (*model.Book, error)
	Delete(ctx context.Context, id int) error
}

// BookService implementa los casos de uso relacionados con libros.
type BookService struct {
	store BookStore
}

// NewBookService crea un BookService respaldado por el store dado.
func NewBookService(s BookStore) *BookService {
	return &BookService{store: s}
}

// ListBooks devuelve una página de libros según los parámetros de búsqueda.
func (s *BookService) ListBooks(ctx context.Context, params model.ListParams) (model.Page[*model.Book], error) {
	params.Normalize()

	books, total, err := s.store.List(ctx, params)
	if err != nil {
		return model.Page[*model.Book]{}, fmt.Errorf("service: listando libros: %w", err)
	}

	return model.NewPage(books, params, total), nil
}

// GetBook obtiene un libro por su ID.
func (s *BookService) GetBook(ctx context.Context, id int) (*model.Book, error) {
	book, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, translateStoreErr(err, "service: obteniendo libro")
	}
	return book, nil
}

// CreateBook valida y crea un nuevo libro.
func (s *BookService) CreateBook(ctx context.Context, book model.Book) (*model.Book, error) {
	book.Sanitize()
	if err := NewValidationError(book.Validate()); err != nil {
		return nil, err
	}

	created, err := s.store.Create(ctx, &book)
	if err != nil {
		return nil, fmt.Errorf("service: creando libro: %w", err)
	}
	return created, nil
}

// UpdateBook valida y actualiza un libro existente.
func (s *BookService) UpdateBook(ctx context.Context, id int, book model.Book) (*model.Book, error) {
	book.Sanitize()
	if err := NewValidationError(book.Validate()); err != nil {
		return nil, err
	}

	updated, err := s.store.Update(ctx, id, &book)
	if err != nil {
		return nil, translateStoreErr(err, "service: actualizando libro")
	}
	return updated, nil
}

// DeleteBook elimina un libro por su ID.
func (s *BookService) DeleteBook(ctx context.Context, id int) error {
	if err := s.store.Delete(ctx, id); err != nil {
		return translateStoreErr(err, "service: eliminando libro")
	}
	return nil
}

func translateStoreErr(err error, context string) error {
	if errors.Is(err, store.ErrNotFound) {
		return ErrNotFound
	}
	return fmt.Errorf("%s: %w", context, err)
}
