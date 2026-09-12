package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/andrew-fe/golang-back/internal/model"
	"github.com/andrew-fe/golang-back/internal/service"
	"github.com/andrew-fe/golang-back/internal/store"
)

// fakeBookStore es una implementación en memoria de service.BookStore
// usada para probar la lógica de negocio sin depender de una base de datos.
type fakeBookStore struct {
	books  map[int]*model.Book
	nextID int
}

func newFakeBookStore() *fakeBookStore {
	return &fakeBookStore{books: make(map[int]*model.Book), nextID: 1}
}

func (f *fakeBookStore) List(_ context.Context, params model.ListParams) ([]*model.Book, int, error) {
	var books []*model.Book
	for _, b := range f.books {
		books = append(books, b)
	}
	return books, len(books), nil
}

func (f *fakeBookStore) GetByID(_ context.Context, id int) (*model.Book, error) {
	b, ok := f.books[id]
	if !ok {
		return nil, store.ErrNotFound
	}
	return b, nil
}

func (f *fakeBookStore) Create(_ context.Context, book *model.Book) (*model.Book, error) {
	book.ID = f.nextID
	f.nextID++
	f.books[book.ID] = book
	return book, nil
}

func (f *fakeBookStore) Update(_ context.Context, id int, book *model.Book) (*model.Book, error) {
	if _, ok := f.books[id]; !ok {
		return nil, store.ErrNotFound
	}
	book.ID = id
	f.books[id] = book
	return book, nil
}

func (f *fakeBookStore) Delete(_ context.Context, id int) error {
	if _, ok := f.books[id]; !ok {
		return store.ErrNotFound
	}
	delete(f.books, id)
	return nil
}

func TestCreateBook_ValidationError(t *testing.T) {
	svc := service.NewBookService(newFakeBookStore())

	_, err := svc.CreateBook(context.Background(), model.Book{Title: "", Author: ""})

	var validationErr *service.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("se esperaba un *service.ValidationError, se obtuvo: %v", err)
	}
	if _, ok := validationErr.Fields["title"]; !ok {
		t.Errorf("se esperaba un error de validación para el campo 'title'")
	}
	if _, ok := validationErr.Fields["author"]; !ok {
		t.Errorf("se esperaba un error de validación para el campo 'author'")
	}
}

func TestCreateBook_Success(t *testing.T) {
	svc := service.NewBookService(newFakeBookStore())

	book, err := svc.CreateBook(context.Background(), model.Book{
		Title:  "  El Aleph  ",
		Author: "Jorge Luis Borges",
		Year:   1949,
	})
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}
	if book.ID == 0 {
		t.Errorf("se esperaba un ID asignado")
	}
	if book.Title != "El Aleph" {
		t.Errorf("se esperaba que el título fuera saneado, se obtuvo: %q", book.Title)
	}
}

func TestGetBook_NotFound(t *testing.T) {
	svc := service.NewBookService(newFakeBookStore())

	_, err := svc.GetBook(context.Background(), 999)

	if !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("se esperaba service.ErrNotFound, se obtuvo: %v", err)
	}
}

func TestUpdateBook_NotFound(t *testing.T) {
	svc := service.NewBookService(newFakeBookStore())

	_, err := svc.UpdateBook(context.Background(), 999, model.Book{Title: "X", Author: "Y"})

	if !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("se esperaba service.ErrNotFound, se obtuvo: %v", err)
	}
}

func TestDeleteBook_NotFound(t *testing.T) {
	svc := service.NewBookService(newFakeBookStore())

	err := svc.DeleteBook(context.Background(), 999)

	if !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("se esperaba service.ErrNotFound, se obtuvo: %v", err)
	}
}

func TestListBooks_Pagination(t *testing.T) {
	st := newFakeBookStore()
	svc := service.NewBookService(st)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		if _, err := svc.CreateBook(ctx, model.Book{Title: "Libro", Author: "Autor"}); err != nil {
			t.Fatalf("no se esperaba error creando datos de prueba: %v", err)
		}
	}

	page, err := svc.ListBooks(ctx, model.ListParams{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}
	if page.TotalItems != 3 {
		t.Errorf("se esperaban 3 items, se obtuvieron %d", page.TotalItems)
	}
}
