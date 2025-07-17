package dto

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// ParsePagination extracts pagination parameters from echo.Context
func ParsePagination(c echo.Context) PaginationParams {
	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("page_size")
	page := 1
	pageSize := 10
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
		pageSize = ps
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return PaginationParams{Page: page, PageSize: pageSize}
}

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

type PaginationParams struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

func (p *PaginationParams) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func NewPaginatedResponse[T any](data []T, page, pageSize int, totalItems int64) *PaginatedResponse[T] {
	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))

	return &PaginatedResponse[T]{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
