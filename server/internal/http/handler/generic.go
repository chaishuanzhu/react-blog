package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"blog-server/internal/apperr"
	"blog-server/internal/http/response"
	"blog-server/internal/model"
)

func idParam(c *gin.Context) (uint64, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, apperr.BadRequest("id must be a positive integer")
	}
	return id, nil
}

func bindJSON[T any](c *gin.Context) (T, bool) {
	var in T
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, bindError(err))
		return in, false
	}
	return in, true
}

func bindError(err error) *apperr.Error {
	if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
		return apperr.ErrTooLarge
	}
	return apperr.BadRequest("invalid JSON body")
}

func CreateHandler[T any](fn func(context.Context, T) (uint64, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		in, ok := bindJSON[T](c)
		if !ok {
			return
		}
		id, err := fn(c.Request.Context(), in)
		if err != nil {
			response.Fail(c, err)
			return
		}
		response.Created(c, model.Created{ID: id})
	}
}

func UpdateHandler[T any](fn func(context.Context, uint64, T) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := idParam(c)
		if err != nil {
			response.Fail(c, err)
			return
		}
		in, ok := bindJSON[T](c)
		if !ok {
			return
		}
		if err := fn(c.Request.Context(), id, in); err != nil {
			response.Fail(c, err)
			return
		}
		response.NoContent(c)
	}
}

func DeleteHandler(fn func(context.Context, uint64) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := idParam(c)
		if err != nil {
			response.Fail(c, err)
			return
		}
		if err := fn(c.Request.Context(), id); err != nil {
			response.Fail(c, err)
			return
		}
		response.NoContent(c)
	}
}
