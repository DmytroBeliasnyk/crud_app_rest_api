package handlers

import (
	"net/http"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Create(ctx *gin.Context) {
	var input dto.ProjectDTO
	if err := ctx.BindJSON(&input); err != nil {
		newErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.ProjectService.Create(input)
	if err != nil {
		newErrResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetById(ctx *gin.Context) {

}

func (h *Handler) GetAll(ctx *gin.Context) {

}

func (h *Handler) UpdateById(ctx *gin.Context) {

}

func (h *Handler) DeleteById(ctx *gin.Context) {

}
