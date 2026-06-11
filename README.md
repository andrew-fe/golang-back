# Backend Go en producción - Estructura de ejemplo completa

Este proyecto es un ejemplo práctico y **totalmente funcional** de cómo organizar un backend en Go siguiendo buenas prácticas para entornos de producción.  
Muestra una arquitectura limpia (clean architecture) con separación de responsabilidades.

## Estructura del proyecto
main/
├── go.mod
├── go.sum
├── main.go
├── README.md
└── internal/
    ├── model/
    │ └── book.go
    ├── service/
    │ └── book_service.go
    ├── store/
    │ └── book_store.go
    └── transport/
      └── book_handler.go


## Contenido de cada archivo

### `go.mod`

```go
module main

go 1.21

require github.com/lib/pq v1.10.9