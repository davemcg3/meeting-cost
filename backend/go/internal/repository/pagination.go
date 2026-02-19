package repository

// Pagination is a common pagination configuration used by repositories.
type Pagination struct {
	Page     int
	PageSize int
	SortBy   string
	SortDir  string // "asc" or "desc"
}

func (p Pagination) Offset() int {
	if p.Page <= 1 {
		return 0
	}
	return (p.Page - 1) * p.PageSize
}

func (p Pagination) Limit() int {
	return p.PageSize
}

