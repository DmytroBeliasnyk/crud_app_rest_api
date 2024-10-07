package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	id, err := h.service.AuthService.SignUp(input)
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

	token, err := h.service.AuthService.SignIn(input)
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

	id, err := h.parseToken(auth)
	if err != nil {
		newErrResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	fmt.Printf("user %d authorized\n", id)
	ctx.Set("user_id", id)
}

func (h *Handler) parseToken(header []string) (int64, error) {
	token, err := jwt.Parse(header[1], func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(h.cfg.Auth.Signature), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token")
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return 0, errors.New("invalid token")
	}

	id, err := strconv.Atoi(sub)
	if err != nil {
		return 0, errors.New("invalid token")
	}

	return int64(id), nil
}
