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

	token, err := h.service.UserService.SignIn(input)
	if err != nil {
		newErrResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"Bearer": token,
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

	id, err := h.auth.ParseToken(auth)
	if err != nil {
		newErrResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.Set("user_id", id)
}
