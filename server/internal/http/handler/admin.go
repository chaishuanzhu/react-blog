package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"blog-server/internal/apperr"
	"blog-server/internal/http/response"
	"blog-server/internal/model"
	"blog-server/internal/service"
)

const defaultAdminPageSize = 12

type Admin struct {
	Svc     *service.Admin
	Uploads *service.Uploads
}

func (h *Admin) Stats(c *gin.Context) {
	stats, err := h.Svc.Stats(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, stats)
}

func (h *Admin) ListArticles(c *gin.Context) {
	page, pageSize, err := pagination(c, defaultAdminPageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	f := model.AdminArticleFilter{Page: page, PageSize: pageSize, Status: c.Query("status")}
	if f.Keyword, err = trimmedQuery(c, "keyword", 64); err != nil {
		response.Fail(c, err)
		return
	}
	if f.CategoryID, err = optionalID(c, "categoryId"); err != nil {
		response.Fail(c, err)
		return
	}
	if f.TagID, err = optionalID(c, "tagId"); err != nil {
		response.Fail(c, err)
		return
	}

	items, total, err := h.Svc.ListArticles(c.Request.Context(), f)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, response.Page[model.AdminArticle]{Items: items, Total: total, Page: page, PageSize: pageSize})
}

func (h *Admin) GetArticle(c *gin.Context) {
	id, err := idParam(c)
	if err != nil {
		response.Fail(c, err)
		return
	}
	article, err := h.Svc.GetArticle(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, article)
}

func (h *Admin) ListComments(c *gin.Context) {
	page, pageSize, err := pagination(c, defaultAdminPageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	items, total, err := h.Svc.ListComments(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, response.Page[model.AdminComment]{Items: items, Total: total, Page: page, PageSize: pageSize})
}

func (h *Admin) UpdatePage(c *gin.Context) {
	in, ok := bindJSON[model.ContentInput](c)
	if !ok {
		return
	}
	if err := h.Svc.UpdatePage(c.Request.Context(), c.Param("key"), in); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContent(c)
}

func (h *Admin) UpdateNotice(c *gin.Context) {
	in, ok := bindJSON[model.NoticeInput](c)
	if !ok {
		return
	}
	if err := h.Svc.UpdateNotice(c.Request.Context(), in); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContent(c)
}

func (h *Admin) UploadTicket(c *gin.Context) {
	in, ok := bindJSON[struct {
		Filename string `json:"filename"`
		Size     int64  `json:"size"`
	}](c)
	if !ok {
		return
	}
	ticket, err := h.Uploads.Ticket(c.Request.Context(), in.Filename, in.Size, time.Now())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, ticket)
}

func optionalID(c *gin.Context, key string) (uint64, error) {
	raw := c.Query(key)
	if raw == "" {
		return 0, nil
	}
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return 0, apperr.BadRequest(key + " must be a positive integer")
	}
	return id, nil
}
