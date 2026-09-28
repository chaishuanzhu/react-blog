package handler

import (
	"github.com/gin-gonic/gin"

	"blog-server/internal/http/response"
	"blog-server/internal/model"
	"blog-server/internal/service"
)

const (
	defaultArticlePageSize = 10
	defaultMomentPageSize  = 20
)

type Public struct {
	Svc *service.Public
}

func (h *Public) ListArticles(c *gin.Context) {
	page, pageSize, err := pagination(c, defaultArticlePageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	f := model.ArticleFilter{Page: page, PageSize: pageSize}
	for key, dst := range map[string]*string{"keyword": &f.Keyword, "category": &f.Category, "tag": &f.Tag} {
		if *dst, err = trimmedQuery(c, key, 64); err != nil {
			response.Fail(c, err)
			return
		}
	}

	items, total, err := h.Svc.ListArticles(c.Request.Context(), f)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, response.Page[model.ArticleSummary]{Items: items, Total: total, Page: page, PageSize: pageSize})
}

func (h *Public) GetArticle(c *gin.Context) {
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

func (h *Public) ListCategories(c *gin.Context) {
	list, err := h.Svc.ListCategories(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}

func (h *Public) ListTags(c *gin.Context) {
	tags, err := h.Svc.ListTags(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"items": tags})
}

func (h *Public) ListMoments(c *gin.Context) {
	page, pageSize, err := pagination(c, defaultMomentPageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	items, total, err := h.Svc.ListMoments(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, response.Page[model.Moment]{Items: items, Total: total, Page: page, PageSize: pageSize})
}

func (h *Public) ListFriendLinks(c *gin.Context) {
	items, err := h.Svc.ListFriendLinks(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"items": items})
}

func (h *Public) ListChangelogs(c *gin.Context) {
	items, err := h.Svc.ListChangelogs(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"items": items})
}

func (h *Public) ListProjects(c *gin.Context) {
	items, err := h.Svc.ListProjects(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"items": items})
}

func (h *Public) GetPage(c *gin.Context) {
	page, err := h.Svc.GetPage(c.Request.Context(), c.Param("key"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, page)
}

func (h *Public) GetSite(c *gin.Context) {
	site, err := h.Svc.GetSite(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, site)
}

func (h *Public) RecordView(c *gin.Context) {
	n, err := h.Svc.RecordView(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"viewCount": n})
}
