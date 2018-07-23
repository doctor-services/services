package paging

import (
	"reflect"
	"testing"
)

func TestGeneratePaginationInfo(t *testing.T) {
	expectedPaginationInfo := Paginator{
		Total:           35,
		TotalPage:       4,
		HasNextPage:     true,
		HasPreviousPage: false,
		NextPage:        2,
	}

	actualPaginationInfo := NewPaginator(35, 10, 1)

	if !reflect.DeepEqual(expectedPaginationInfo, actualPaginationInfo) {
		t.Fatalf("Expected %v but got %v", expectedPaginationInfo, actualPaginationInfo)
	}
	expectedPaginationInfo = Paginator{
		Total:           35,
		TotalPage:       4,
		HasNextPage:     true,
		HasPreviousPage: true,
		NextPage:        3,
		PreviousPage:    1,
	}

	actualPaginationInfo = NewPaginator(35, 10, 2)

	if !reflect.DeepEqual(expectedPaginationInfo, actualPaginationInfo) {
		t.Fatalf("Expected %v but got %v", expectedPaginationInfo, actualPaginationInfo)
	}
}
