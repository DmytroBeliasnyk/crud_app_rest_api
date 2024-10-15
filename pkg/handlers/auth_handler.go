package handlers

import (
	"net/http"
	"strings"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/gin-gonic/gin"
)

// signUp godoc
//
//	@Summary	signUp
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		SignUpDTO	body		dto.SignUpDTO	true	"user details"
//	@Success	201			{integer}	integer			user_id
//	@Failure	400			{object}	errResponse
//	@Failure	500			{object}	errResponse
//	@Failure	default		{object}	errResponse
//	@Router		/auth/sign-up [post]
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

	id, err := h.service.AuthService.SignUp(input)
	if err != nil {
		newErrResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

// signIn godoc
//
//	@Summary	signIn
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		SignInDTO	body		dto.SignInDTO	true	"user details"
//	@Header		200			{string}	Set-Cookie		"set new refresh token"
//	@Success	200			{string}	string			jwt
//	@Failure	400			{object}	errResponse
//	@Failure	401			{object}	errResponse
//	@Failure	default		{object}	errResponse
//	@Router		/auth/sign-in [post]
func (h *Handler) signIn(ctx *gin.Context) {
	var input dto.SignInDTO
	if err := ctx.BindJSON(&input); err != nil {
		newErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := h.service.AuthService.SignIn(input)
	if err != nil {
		newErrResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	jt, rt, err := h.service.AuthService.GenerateTokens(userId)
	if err != nil {
		newErrResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SetCookie(h.cfg.name, rt, h.cfg.age, h.cfg.path, h.cfg.domain, h.cfg.secure, h.cfg.httpOnly)
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"Bearer": jt,
	})
}

// refresh godoc
//
//	@Summary		refresh
//	@Description	refreshing jwt
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			refresh-token	header		string		true	"refresh token from cookie"
//	@Header			200				{string}	Set-Cookie	"Set new refresh token"
//	@Success		200				{string}	string		jwt
//	@Failure		400				{object}	errResponse
//	@Failure		default			{object}	errResponse
//	@Router			/auth/refresh [get]
func (h *Handler) refresh(ctx *gin.Context) {
	rt, err := ctx.Cookie("refresh-token")
	if err != nil {
		newErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	jt, rt, err := h.service.AuthService.UpdateTokens(rt)
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

	id, err := h.service.AuthService.ParseToken(auth[1])
	if err != nil {
		newErrResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.Set("user_id", id)
}
