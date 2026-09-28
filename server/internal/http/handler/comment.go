package handler

import (
	"github.com/gin-gonic/gin"

	"blog-server/internal/http/middleware"
	"blog-server/internal/http/response"
	"blog-server/internal/model"
	"blog-server/internal/service"
)

const defaultCommentPageSize = 10

type Comment struct {
	Svc  *service.Comments
	Auth *service.Auth
}

func (h *Comment) List(c *gin.Context) {
	page, pageSize, err := pagination(c, defaultCommentPageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	articleID, err := optionalIDQuery(c, "articleId")
	if err != nil {
		response.Fail(c, err)
		return
	}
	items, total, err := h.Svc.List(c.Request.Context(), articleID, page, pageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, response.Page[model.CommentThread]{Items: items, Total: total, Page: page, PageSize: pageSize})
}

func (h *Comment) Create(c *gin.Context) {
	var in model.CommentInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, bindError(err))
		return
	}
	author := model.CommentAuthor{IP: c.ClientIP(), UserAgent: c.Request.UserAgent()}
	if userID, ok := middleware.UserID(c); ok {
		user, err := h.Auth.GetUser(c.Request.Context(), userID)
		if err != nil {
			response.Fail(c, err)
			return
		}
		author.IsAdmin = true
		author.Nickname, author.Email, author.Website, author.Avatar = user.Nickname, user.Email, user.Website, user.Avatar
	}

	comment, err := h.Svc.Create(c.Request.Context(), in, author)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, comment)
}
