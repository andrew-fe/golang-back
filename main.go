package main

import (
	"database/sql"
	"fmt"
	"log"
	"main/internal/service"
	"main/internal/store"
	"main/internal/transport"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	// Conectar a PostgreSQL
	db, err := sql.Open("postgres", "user=postgres password=Andres1234 dbname=Golang sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Crear el table si no existe
	q := `
		CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL
		)
	`
	if _, err := db.Exec(q); err != nil {
		log.Fatal(err.Error())
	}

	// inyectar nuestra dependencias
	bookStore := store.New(db)
	bookService := service.New(bookStore)
	bookHandler := transport.New(bookService)

	// Configurar las rutas
	http.HandleFunc("/books", bookHandler.HandleBooks)
	http.HandleFunc("/books/", bookHandler.HandleBookByID)

	fmt.Println("🚀 Servidor ejecutandos en http://localhost:8080")
	fmt.Println("📚 main Endpoints:")
	fmt.Println("   GET    /books 	   	- Obtener todos los libros")
	fmt.Println("   POST   /books 	   	- Crear un nuevo libro")
	fmt.Println("   GET    /books/{id} 		- Obtener un libro por ID")
	fmt.Println("   PUT    /books/{id} 		- Actualizar un libro por ID")
	fmt.Println("   DELETE /books/{id} 		- Eliminar un libro por ID")

	// Iniciar el servidor
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
