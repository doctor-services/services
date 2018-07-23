package paging

import "math"

// Paginator helps to calculate paging infor
type Paginator struct {
	Total           int  `json:"total"`
	TotalPage       int  `json:"totalPage"`
	HasNextPage     bool `json:"hasNextPage"`
	HasPreviousPage bool `json:"hasPreviousPage"`
	NextPage        int  `json:"nextPage,omitempty"`
	PreviousPage    int  `json:"previousPage,omitempty"`
}

// NewPaginator creates the pagination information
func NewPaginator(total int, limit int, currentPage int) Paginator {
	pageNums := int(math.Ceil(float64(total) / float64(limit)))
	result := Paginator{
		Total:     total,
		TotalPage: pageNums,
	}
	if pageNums > currentPage {
		result.HasNextPage = true
		result.NextPage = currentPage + 1
	}
	if currentPage > 1 {
		result.HasPreviousPage = true
		result.PreviousPage = currentPage - 1
	}
	return result
}
