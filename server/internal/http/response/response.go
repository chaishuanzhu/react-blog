package response

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"blog-server/internal/apperr"
)

type Page[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, data)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Fail(c *gin.Context, err error) {
	var apiErr *apperr.Error
	if !errors.As(err, &apiErr) {
		slog.ErrorContext(c.Request.Context(), "internal error", "err", err, "path", c.FullPath())
		apiErr = &apperr.Error{Status: http.StatusInternalServerError, Code: "INTERNAL", Message: "internal server error"}
	}
	c.AbortWithStatusJSON(apiErr.Status, gin.H{"error": apiErr})
}
