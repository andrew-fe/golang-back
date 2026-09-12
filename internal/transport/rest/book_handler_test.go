package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/andrew-fe/golang-back/internal/model"
	"github.com/andrew-fe/golang-back/internal/service"
	"github.com/andrew-fe/golang-back/internal/transport/rest"
)

// fakeBookService implementa rest.BookService en memoria para probar los
// handlers HTTP de forma aislada, sin base de datos.
type fakeBookService struct {
	books map[int]*model.Book
}

func newFakeBookService() *fakeBookService {
	return &fakeBookService{books: map[int]*model.Book{
		1: {ID: 1, Title: "Cien años de soledad", Author: "Gabriel García Márquez"},
	}}
}

func (f *fakeBookService) ListBooks(_ context.Context, params model.ListParams) (model.Page[*model.Book], error) {
	var books []*model.Book
	for _, b := range f.books {
		books = append(books, b)
	}
	return model.NewPage(books, params, len(books)), nil
}

func (f *fakeBookService) GetBook(_ context.Context, id int) (*model.Book, error) {
	b, ok := f.books[id]
	if !ok {
		return nil, service.ErrNotFound
	}
	return b, nil
}

func (f *fakeBookService) CreateBook(_ context.Context, book model.Book) (*model.Book, error) {
	if fields := (&book).Validate(); len(fields) > 0 {
		return nil, &service.ValidationError{Fields: fields}
	}
	book.ID = len(f.books) + 1
	f.books[book.ID] = &book
	return &book, nil
}

func (f *fakeBookService) UpdateBook(_ context.Context, id int, book model.Book) (*model.Book, error) {
	if _, ok := f.books[id]; !ok {
		return nil, service.ErrNotFound
	}
	book.ID = id
	f.books[id] = &book
	return &book, nil
}

func (f *fakeBookService) DeleteBook(_ context.Context, id int) error {
	if _, ok := f.books[id]; !ok {
		return service.ErrNotFound
	}
	delete(f.books, id)
	return nil
}

func newTestRouter(svc rest.BookService) http.Handler {
	r := chi.NewRouter()
	h := rest.NewBookHandler(svc)
	r.Route("/api/v1/books", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
	return r
}

func TestBookHandler_Get_NotFound(t *testing.T) {
	router := newTestRouter(newFakeBookService())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/books/999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba status 404, se obtuvo %d", rec.Code)
	}
}

func TestBookHandler_Get_Success(t *testing.T) {
	router := newTestRouter(newFakeBookService())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/books/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba status 200, se obtuvo %d", rec.Code)
	}

	var got model.Book
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("se esperaba el libro con ID 1, se obtuvo %d", got.ID)
	}
}

func TestBookHandler_Create_ValidationError(t *testing.T) {
	router := newTestRouter(newFakeBookService())

	body, _ := json.Marshal(model.Book{Title: "", Author: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/books/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("se esperaba status 422, se obtuvo %d", rec.Code)
	}
}

func TestBookHandler_Create_Success(t *testing.T) {
	router := newTestRouter(newFakeBookService())

	body, _ := json.Marshal(model.Book{Title: "1984", Author: "George Orwell", Year: 1949})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/books/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba status 201, se obtuvo %d: %s", rec.Code, rec.Body.String())
	}
}

func TestBookHandler_Delete_NotFound(t *testing.T) {
	router := newTestRouter(newFakeBookService())

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/books/999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba status 404, se obtuvo %d", rec.Code)
	}
}
