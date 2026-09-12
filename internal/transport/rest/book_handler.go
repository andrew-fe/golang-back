package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/andrew-fe/golang-back/internal/model"
	"github.com/andrew-fe/golang-back/internal/service"
)

// BookService es el contrato de negocio requerido por BookHandler.
type BookService interface {
	ListBooks(ctx context.Context, params model.ListParams) (model.Page[*model.Book], error)
	GetBook(ctx context.Context, id int) (*model.Book, error)
	CreateBook(ctx context.Context, book model.Book) (*model.Book, error)
	UpdateBook(ctx context.Context, id int, book model.Book) (*model.Book, error)
	DeleteBook(ctx context.Context, id int) error
}

// BookHandler expone los endpoints HTTP del recurso "books".
type BookHandler struct {
	service BookService
}

// NewBookHandler crea un BookHandler respaldado por el servicio dado.
func NewBookHandler(s BookService) *BookHandler {
	return &BookHandler{service: s}
}

// List maneja GET /books, soportando búsqueda (?q=) y paginación (?page=, ?page_size=).
func (h *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	params := model.ListParams{
		Search:   r.URL.Query().Get("q"),
		Page:     atoiOrDefault(r.URL.Query().Get("page"), 1),
		PageSize: atoiOrDefault(r.URL.Query().Get("page_size"), 20),
	}

	page, err := h.service.ListBooks(r.Context(), params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "no se pudo obtener la lista de libros")
		return
	}

	writeJSON(w, http.StatusOK, page)
}

// Get maneja GET /books/{id}.
func (h *BookHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := idFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "identificador de libro no válido")
		return
	}

	book, err := h.service.GetBook(r.Context(), id)
	if err != nil {
		handleServiceError(w, err, "no se pudo obtener el libro")
		return
	}

	writeJSON(w, http.StatusOK, book)
}

// Create maneja POST /books.
func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var book model.Book
	if err := decodeJSON(r, &book); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo de la petición inválido")
		return
	}

	created, err := h.service.CreateBook(r.Context(), book)
	if err != nil {
		handleServiceError(w, err, "no se pudo crear el libro")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// Update maneja PUT /books/{id}.
func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := idFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "identificador de libro no válido")
		return
	}

	var book model.Book
	if err := decodeJSON(r, &book); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo de la petición inválido")
		return
	}

	updated, err := h.service.UpdateBook(r.Context(), id, book)
	if err != nil {
		handleServiceError(w, err, "no se pudo actualizar el libro")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// Delete maneja DELETE /books/{id}.
func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := idFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "identificador de libro no válido")
		return
	}

	if err := h.service.DeleteBook(r.Context(), id); err != nil {
		handleServiceError(w, err, "no se pudo eliminar el libro")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func idFromRequest(r *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(r, "id"))
}

func decodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func atoiOrDefault(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return n
}

func handleServiceError(w http.ResponseWriter, err error, fallbackMessage string) {
	var validationErr *service.ValidationError

	switch {
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, "libro no encontrado")
	case errors.As(err, &validationErr):
		writeValidationError(w, validationErr.Fields)
	default:
		slog.Error(fallbackMessage, "error", err)
		writeError(w, http.StatusInternalServerError, fallbackMessage)
	}
}
