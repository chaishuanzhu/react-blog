package handler

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"blog-server/internal/apperr"
)

const maxPageSize = 100

func pagination(c *gin.Context, defaultSize int) (page, pageSize int, err error) {
	page, err = positiveInt(c.Query("page"), 1)
	if err != nil {
		return 0, 0, apperr.BadRequest("page must be a positive integer")
	}
	pageSize, err = positiveInt(c.Query("pageSize"), defaultSize)
	if err != nil || pageSize > maxPageSize {
		return 0, 0, apperr.BadRequest("pageSize must be between 1 and 100")
	}
	return page, pageSize, nil
}

func positiveInt(raw string, fallback int) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 0, strconv.ErrRange
	}
	return n, nil
}

func trimmedQuery(c *gin.Context, key string, maxLen int) (string, error) {
	v := strings.TrimSpace(c.Query(key))
	if len([]rune(v)) > maxLen {
		return "", apperr.BadRequest(key + " is too long")
	}
	return v, nil
}

// optionalIDQuery returns 0 when the query parameter is absent.
func optionalIDQuery(c *gin.Context, key string) (uint64, error) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return 0, nil
	}
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return 0, apperr.BadRequest(key + " must be a positive integer")
	}
	return id, nil
}
