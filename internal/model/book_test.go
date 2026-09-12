package model_test

import (
	"testing"

	"github.com/andrew-fe/golang-back/internal/model"
)

func TestBook_Validate(t *testing.T) {
	tests := []struct {
		name       string
		book       model.Book
		wantFields []string
	}{
		{
			name:       "libro válido",
			book:       model.Book{Title: "Rayuela", Author: "Julio Cortázar", Year: 1963},
			wantFields: nil,
		},
		{
			name:       "sin título ni autor",
			book:       model.Book{Title: "", Author: ""},
			wantFields: []string{"title", "author"},
		},
		{
			name:       "año futuro no válido",
			book:       model.Book{Title: "X", Author: "Y", Year: 9999},
			wantFields: []string{"year"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.book.Validate()
			if len(errs) != len(tt.wantFields) {
				t.Fatalf("se esperaban errores en %v, se obtuvo %v", tt.wantFields, errs)
			}
			for _, field := range tt.wantFields {
				if _, ok := errs[field]; !ok {
					t.Errorf("se esperaba un error en el campo %q", field)
				}
			}
		})
	}
}
