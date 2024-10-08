package handlers

import (
	"net/http"
	"strings"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(ctx *gin.Context) {
	var input dto.SignUpDTO
	if err := ctx.BindJSON(&input); err != nil {
		newErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		newErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.UserService.SignUp(input)
	if err != nil {
		newErrResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) signIn(ctx *gin.Context) {
	var input dto.SignInDTO
	if err := ctx.BindJSON(&input); err != nil {
		newErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	jt, rt, err := h.service.UserService.SignIn(input)
	if err != nil {
		newErrResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.SetCookie(h.cfg.name, rt, h.cfg.age, h.cfg.path, h.cfg.domain, h.cfg.secure, h.cfg.httpOnly)
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"Bearer": jt,
	})
}

func (h *Handler) refresh(ctx *gin.Context) {
	rt, err := ctx.Cookie("refresh-token")
	if err != nil {
		newErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	jt, rt, err := h.auth.UpdateTokens(rt)
	if err != nil {
		newErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.SetCookie(h.cfg.name, rt, h.cfg.age, h.cfg.path, h.cfg.domain, h.cfg.secure, h.cfg.httpOnly)
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"Bearer": jt,
	})
}

func (h *Handler) middlewareAuth(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		newErrResponse(ctx, http.StatusBadRequest, "empty authorization header")
		return
	}

	auth := strings.Split(header, " ")
	if len(auth) != 2 || auth[0] != "Bearer" {
		newErrResponse(ctx, http.StatusBadRequest, "invalid authorization header")
		return
	}

	id, err := h.auth.ParseToken(auth[1])
	if err != nil {
		newErrResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.Set("user_id", id)
}
