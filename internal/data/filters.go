package data

import (
	"greenlight/internal/validator"
	"math"
	"strings"
)

// Filters provide a means of filtering SQL queries
type Filters struct {
	Page         int // Consider changing => PageIndex
	PageSize     int
	Sort         string
	SortSafeList []string
}

// Metadata holds Additional data provided to client indicating Filters
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter" + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

// offset _ sets the SQL page (index).
//We need to Subtract 1 from the PageIndex Offset to account for the fact that
// in SQL OFFSET inital value = 0 i.e. it is Zero-Indexed.
// Otherwise par example. SELECT * FROM movies LIMIT 1 OFFSET 1 returns entry
// with id=2. Hence we are required to offset index by 1.
// https://wellingguzman.com/notes/pagination-with-mysql
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// ValidateFilters validates a given page filter
func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize < 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(validator.In(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}
