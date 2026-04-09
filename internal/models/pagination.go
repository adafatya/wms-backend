package models

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page      int  `json:"page"`
	Limit     int  `json:"limit"`
	PrevPage  *int `json:"prev_page"`
	NextPage  *int `json:"next_page"`
	TotalPage int  `json:"total_page"`
}

func NewPagination(page, limit int, totalData int64) *Pagination {
	if limit <= 0 {
		limit = 10
	}
	totalPage := 0
	if totalData > 0 {
		totalPage = (int(totalData) + limit - 1) / limit
	}

	var prevPage *int
	if page > 1 {
		p := page - 1
		prevPage = &p
	}

	var nextPage *int
	if page < totalPage {
		n := page + 1
		nextPage = &n
	}

	return &Pagination{
		Page:      page,
		Limit:     limit,
		PrevPage:  prevPage,
		NextPage:  nextPage,
		TotalPage: totalPage,
	}
}

// ParsePaginationQuery extracts and validates page and limit from query parameters
func ParsePaginationQuery(c *gin.Context) (int, int, error) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, 0, errors.New("invalid page parameter: must be a positive integer greater than 0")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return 0, 0, errors.New("invalid limit parameter: must be a positive integer greater than 0")
	}

	if limit > 100 {
		return 0, 0, fmt.Errorf("limit too large: maximum allowed is 100")
	}

	return page, limit, nil
}
