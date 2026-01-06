package util

import (
	"database/sql"
)

type PaginatedResponse[T any] struct {
	Data  []*T `json:"data"`
	Total int  `json:"total"`
}

func NewPagination(page, pageSize int) (offset, limit int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset = (page - 1) * pageSize
	limit = pageSize
	return
}

// SỬA Ở ĐÂY
func Paginate[T any](
	db *sql.DB,
	query string,
	countQuery string,
	filterArgs []interface{},
	offset int,
	limit int,
	scanRow func(*sql.Rows) (*T, error),
) (*PaginatedResponse[T], error) {

	// 1️Count tổng (CHỈ filter)
	var total int
	if err := db.QueryRow(countQuery, filterArgs...).Scan(&total); err != nil {
		return nil, err
	}

	//  Query data (filter + limit + offset)
	argsWithLimit := append(
		append([]interface{}{}, filterArgs...),
		limit,
		offset,
	)

	rows, err := db.Query(query, argsWithLimit...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*T
	for rows.Next() {
		item, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return &PaginatedResponse[T]{
		Data:  data,
		Total: total,
	}, nil
}
