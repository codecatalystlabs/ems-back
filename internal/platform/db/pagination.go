package db

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Pagination struct {
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	Offset    int               `json:"-"`
	Search    string            `json:"search"`
	SortBy    string            `json:"sort_by"`
	SortOrder string            `json:"sort_order"`
	Filters   map[string]string `json:"filters"`
}

type PageMeta struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	Total      int64  `json:"total"`
	TotalPages int    `json:"total_pages"`
	Search     string `json:"search,omitempty"`
	SortBy     string `json:"sort_by,omitempty"`
	SortOrder  string `json:"sort_order,omitempty"`
}

type PageResult[T any] struct {
	Items []T      `json:"items"`
	Meta  PageMeta `json:"meta"`
}

func ParsePagination(values url.Values, allowedSortFields map[string]string, allowedFilters map[string]struct{}) Pagination {
	page := parsePositiveInt(values.Get("page"), 1)
	pageSize := parsePositiveInt(values.Get("page_size"), 20)
	if pageSize > 100 {
		pageSize = 100
	}

	sortBy := strings.TrimSpace(values.Get("sort_by"))
	if sortBy == "" {
		sortBy = "created_at"
	}
	if _, ok := allowedSortFields[sortBy]; !ok {
		sortBy = "created_at"
	}

	sortOrder := strings.ToUpper(strings.TrimSpace(values.Get("sort_order")))
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	filters := make(map[string]string)
	for key, vals := range values {
		if !strings.HasPrefix(key, "filter[") || !strings.HasSuffix(key, "]") {
			continue
		}
		name := strings.TrimSuffix(strings.TrimPrefix(key, "filter["), "]")
		if _, ok := allowedFilters[name]; !ok {
			continue
		}
		if len(vals) == 0 {
			continue
		}
		v := strings.TrimSpace(vals[0])
		if v != "" {
			filters[name] = v
		}
	}

	return Pagination{
		Page:      page,
		PageSize:  pageSize,
		Offset:    (page - 1) * pageSize,
		Search:    strings.TrimSpace(values.Get("search")),
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Filters:   filters,
	}
}

func BuildOrderBy(p Pagination, allowedSortFields map[string]string) string {
	col, ok := allowedSortFields[p.SortBy]
	if !ok {
		col = allowedSortFields["created_at"]
	}
	order := "DESC"
	if p.SortOrder == "ASC" {
		order = "ASC"
	}
	return fmt.Sprintf("ORDER BY %s %s", col, order)
}

func NewPageMeta(p Pagination, total int64) PageMeta {
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(p.PageSize) - 1) / int64(p.PageSize))
	}
	return PageMeta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: totalPages,
		Search:     p.Search,
		SortBy:     p.SortBy,
		SortOrder:  p.SortOrder,
	}
}

func parsePositiveInt(raw string, fallback int) int {
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}
