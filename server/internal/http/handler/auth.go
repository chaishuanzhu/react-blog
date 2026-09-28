package handler

import (
	"github.com/gin-gonic/gin"

	"blog-server/internal/apperr"
	"blog-server/internal/http/middleware"
	"blog-server/internal/http/response"
	"blog-server/internal/model"
	"blog-server/internal/service"
)

type Auth struct {
	Svc *service.Auth
}

func (h *Auth) Login(c *gin.Context) {
	in, ok := bindJSON[struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}](c)
	if !ok {
		return
	}
	res, err := h.Svc.Login(c.Request.Context(), in.Email, in.Password)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Auth) Me(c *gin.Context) {
	id, _ := middleware.UserID(c)
	user, err := h.Svc.GetUser(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, user)
}

func (h *Auth) UpdateProfile(c *gin.Context) {
	in, ok := bindJSON[model.ProfileInput](c)
	if !ok {
		return
	}
	id, _ := middleware.UserID(c)
	user, err := h.Svc.UpdateProfile(c.Request.Context(), id, in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, user)
}

func (h *Auth) ChangePassword(c *gin.Context) {
	in, ok := bindJSON[struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}](c)
	if !ok {
		return
	}
	if in.CurrentPassword == "" || in.NewPassword == "" {
		response.Fail(c, apperr.BadRequest("currentPassword and newPassword are required"))
		return
	}
	id, _ := middleware.UserID(c)
	res, err := h.Svc.ChangePassword(c.Request.Context(), id, in.CurrentPassword, in.NewPassword)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}
