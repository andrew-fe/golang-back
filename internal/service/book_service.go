package service

import (
	"errors"
	"main/internal/model"
	"main/internal/store"
)

type Service struct {
	store store.Store
}

func New(s store.Store) *Service {
	return &Service{
		store: s,
	}
}

func (s *Service) GetAllBooks() ([]*model.Book, error) {

	books, err := s.store.GetAll()
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (s *Service) GetBookByID(id int) (*model.Book, error) {
	book, err := s.store.GetByID(id)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *Service) CreateBook(book model.Book) (*model.Book, error) {

	if book.Title == "" || book.Author == "" {
		return nil, errors.New("Se requieren el título y el autor del libro")
	}

	bookCreated, err := s.store.Create(&book)
	if err != nil {
		return nil, err
	}

	return bookCreated, nil
}

func (s *Service) UpdateBook(id int, book model.Book) (*model.Book, error) {

	if book.Title == "" || book.Author == "" {
		return nil, errors.New("Se requieren el título y el autor del libro")
	}

	updatedBook, err := s.store.Update(id, &book)
	if err != nil {
		return nil, err
	}
	return updatedBook, nil
}

func (s *Service) DeleteBook(id int) error {
	err := s.store.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
