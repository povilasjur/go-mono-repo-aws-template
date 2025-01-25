package common

import (
	"math"
	"strconv"
)

type SortType string

const (
	SortAsc  SortType = "asc"
	SortDesc SortType = "desc"
)

// swagger:parameters getAll
type PageFilter struct {
	// Page size of the results
	//
	// in: query
	// required: false
	// example: 20
	PageSize int64 `json:"pageSize"`

	// Page number to retrieve
	//
	// in: query
	// required: false
	// example: 1
	Page int64 `json:"page"`

	// Sort by a specific field
	//
	// in: query
	// required: false
	// example: name
	SortField string `json:"sortField"`

	// Sort order, either asc or desc
	//
	// in: query
	// required: false
	// enum: ["asc", "desc"]
	// default: "asc"
	SortType SortType `json:"sortType"`
}

type Paginated[D any] struct {
	Data       []*D            `json:"data"`
	Pagination *PaginationData `json:"page"`
}

func NewPaginated[D any](data []*D, total int64, pageSize int64, page int64) *Paginated[D] {
	return &Paginated[D]{
		Data:       data,
		Pagination: NewPaginationData(total, pageSize, page),
	}
}

type PaginationData struct {
	Total     int64 `json:"total"`
	Number    int64 `json:"number"`
	PageSize  int64 `json:"pageSize"`
	TotalPage int64 `json:"totalPage"`
}

func NewPaginationData(total int64, pageSize int64, pageNumber int64) *PaginationData {
	return &PaginationData{
		Total:     total,
		Number:    pageNumber,
		PageSize:  pageSize,
		TotalPage: int64(math.Ceil(float64(total) / float64(pageSize))),
	}
}

func ParsePageFilter(queryParams map[string]string) *PageFilter {
	pageSizeString := queryParams["pageSize"]
	pageSize := int64(0)
	if pageSizeString != "" {
		pageSize, _ = strconv.ParseInt(pageSizeString, 10, 64)
	}

	pageString := queryParams["page"]
	page := int64(0)
	if pageString != "" {
		page, _ = strconv.ParseInt(pageString, 10, 64)
	}

	sort := queryParams["sort"]
	if sort == "" {
		sort = "_id"
	}

	sortTypeString := queryParams["sortType"]
	sortType := SortAsc
	if sortTypeString == "desc" {
		sortType = SortDesc
	}

	if pageSize < 1 {
		pageSize = 10
	}
	if page < 1 {
		page = 1
	}

	return &PageFilter{
		PageSize:  pageSize,
		Page:      page,
		SortField: sort,
		SortType:  sortType,
	}
}

func (p PageFilter) GetLimit() int64 {
	return p.PageSize * p.Page
}

func (p PageFilter) GetSkip() int64 {
	if p.Page == 0 {
		return 0
	}
	return p.PageSize * (p.Page - 1)
}

func (p PageFilter) GetSortTypeInt() int {
	sortType := 1
	if p.SortType == "desc" {
		sortType = -1
	}
	return sortType
}
