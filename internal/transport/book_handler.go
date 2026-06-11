package transport

import (
	"encoding/json"
	"main/internal/model"
	"main/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type BookHandler struct {
	service *service.Service
}

func New(s *service.Service) *BookHandler {
	return &BookHandler{
		service: s,
	}
}

func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		books, err := h.service.GetAllBooks()
		if err != nil {
			http.Error(w, "Falló al obtener los libros", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)

	case http.MethodPost:
		var book model.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "Datos de libro no válidos", http.StatusBadRequest)
			return
		}

		createdBook, err := h.service.CreateBook(book)
		if err != nil {
			http.Error(w, "Falló al crear el libro", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(createdBook)

	default:
		http.Error(w, "Método no disponible", http.StatusMethodNotAllowed)

	}
}

func (h *BookHandler) HandleBookByID(w http.ResponseWriter, r *http.Request) {
	// Obtener ID de la URL
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de libro no válido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		book, err := h.service.GetBookByID(id)
		if err != nil {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)

	case http.MethodPut:
		var book model.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "Datos de libro no válidos", http.StatusBadRequest)
			return
		}
		updatedBook, err := h.service.UpdateBook(id, book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedBook)

	case http.MethodDelete:
		err := h.service.DeleteBook(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Método no disponible", http.StatusMethodNotAllowed)
	}
}
