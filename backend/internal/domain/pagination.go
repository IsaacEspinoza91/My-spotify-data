package domain

type Pagination struct {
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"` // Aquí irán tus []DTOs
}

// NewPagination asegura que Data nunca sea nil
func NewPagination(data interface{}, total, page, limit int) Pagination {
	if data == nil {
		data = []interface{}{} // Garantiza [] en el JSON
	}

	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return Pagination{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Data:       data,
	}
}
