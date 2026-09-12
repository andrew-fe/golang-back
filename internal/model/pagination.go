package model

// ListParams describe los parámetros de búsqueda y paginación soportados
// por los listados del dominio.
type ListParams struct {
	Search   string
	Page     int
	PageSize int
}

// Normalize aplica valores por defecto y límites razonables a los
// parámetros de paginación.
func (p *ListParams) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// Offset calcula el offset SQL correspondiente a la página actual.
func (p ListParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Page representa una página de resultados junto a metadatos de paginación.
type Page[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// NewPage construye una respuesta paginada calculando el total de páginas.
func NewPage[T any](items []T, params ListParams, totalItems int) Page[T] {
	totalPages := 0
	if params.PageSize > 0 {
		totalPages = (totalItems + params.PageSize - 1) / params.PageSize
	}
	return Page[T]{
		Items:      items,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
