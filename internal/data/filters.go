package data

import (
	"net/url"
	"strconv"
	"strings"

	"PhoceeneAuto/internal/validator"
)

// Filters represents a set of filters for querying data, including pagination and sorting.
type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

// NewPostFilters creates a new instance of Filters based on the query parameters from an HTTP request.
//
// Parameters:
//
//	q - The URL query parameters
//
// Returns:
//
//	*Filters - A pointer to the newly created Filters instance
func NewPostFilters(q url.Values) *Filters {

	// setting the basic post filters
	var filters = &Filters{
		PageSize:     12,
		SortSafelist: []string{"title", "created_at", "updated_at", "id", "-title", "-created_at", "-updated_at", "-id"},
	}

	// getting the page
	if q.Has("page") {
		filters.Page, _ = strconv.Atoi(q.Get("page"))
	} else {
		filters.Page = 1
	}

	// getting the sorting order
	if q.Has("sort") {
		filters.Sort = q.Get("sort")
	} else {
		filters.Sort = "id"
	}

	return filters
}

// sortColumn returns the column name for sorting, removing any prefix indicating direction.
//
// Returns:
//
//	string - The column name for sorting
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

// sortDirection returns the sorting direction ("ASC" or "DESC").
//
// Returns:
//
//	string - The sorting direction
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

// limit returns the number of records to return per page.
//
// Returns:
//
//	int - The number of records per page
func (f Filters) limit() int {
	return f.PageSize
}

// offset calculates the starting point for pagination based on the current page and page size.
//
// Returns:
//
//	int - The offset value
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// ValidateFilters validates the Filters instance using a Validator.
//
// Parameters:
//
//	v - The Validator instance to use for validation
//	f - The Filters instance to validate
func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

// Metadata represents metadata for pagination and filtering results.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// calculateMetadata calculates the pagination metadata based on total records, current page, and page size.
//
// Parameters:
//
//	totalRecords - The total number of records
//	page - The current page number
//	pageSize - The number of records per page
//
// Returns:
//
//	Metadata - The calculated pagination metadata
func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}
