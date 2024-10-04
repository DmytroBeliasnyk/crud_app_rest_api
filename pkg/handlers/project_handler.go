package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/gin-gonic/gin"
)

// Create godoc
//
//	@Summary		Create
//	@Description	create new project
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Param			input	body		dto.ProjectDTO	true	"project info"
//	@Success		201		{integer}	integer			id
//	@Failure		400		{object}	errResponse
//	@Failure		500		{object}	errResponse
//	@Failure		default	{obkect}	errResponse
//	@Router			/api/projects/ [post]
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

	h.cache.Delete("all")

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

// GetById godoc
//
//	@Summary		GetById
//	@Description	get project by id
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Param			id		query		integer	true	"project id"
//	@Success		200		{object}	dto.ProjectDTO
//	@Failure		400		{object}	errResponse
//	@Failure		500		{object}	errResponse
//	@Failure		default	{object}	errResponse
//	@Router			/api/projects [get]
func (h *Handler) GetById(ctx *gin.Context) {
	paramId := ctx.Query("id")

	project, err := h.cache.Get(paramId)
	if err != nil {
		id, err := strconv.Atoi(ctx.Query("id"))
		if err != nil {
			newErrResponse(ctx, http.StatusBadRequest, fmt.Sprintf("%s: message: invalid id param", err))
			return
		}

		project, err = h.service.ProjectService.GetById(int64(id))
		if err != nil {
			newErrResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		h.cache.Set(paramId, project, time.Hour)
	}

	ctx.JSON(http.StatusOK, project)
}

// GetAll godoc
//
//	@Summary		GetAll
//	@Description	get all projects
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Success		200		{array}		dto.ProjectDTO
//	@Failure		500		{object}	errResponse
//	@Failure		default	{object}	errResponse
//	@Router			/api/projects/ [get]
func (h *Handler) GetAll(ctx *gin.Context) {
	projects, err := h.cache.Get("all")
	if err != nil {
		projects, err = h.service.ProjectService.GetAll()
		if err != nil {
			newErrResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		h.cache.Set("all", projects, time.Hour)
	}

	ctx.JSON(http.StatusOK, projects)
}

// UpdateById godoc
//
//	@Summary		UpdateById
//	@Description	update project by id
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Param			id		path		integer			true	"project id"
//	@Param			input	body		dto.ProjectDTO	true	"project info"
//	@Success		200		{object}	statusResponse
//	@Failure		400		{object}	errResponse
//	@Failure		500		{object}	errResponse
//	@Failure		default	{object}	errResponse
//	@Router			/api/projects/{id} [post]
func (h *Handler) UpdateById(ctx *gin.Context) {
	paramId := ctx.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		newErrResponse(ctx, http.StatusBadRequest, fmt.Sprintf("%s: message: invalid id param", err))
		return
	}

	var input dto.UpdateProjectDTO
	if err := ctx.BindJSON(&input); err != nil {
		newErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		newErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.service.ProjectService.UpdateById(int64(id), input); err != nil {
		newErrResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	h.cache.Delete("all")
	h.cache.Delete(paramId)

	ctx.JSON(http.StatusOK, statusResponse{"ok"})
}

// DeleteById godoc
//
//	@Summary		DeleteById
//	@Description	delete project by id
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Param			id		path		integer	true	"project id"
//	@Success		200		{object}	statusResponse
//	@Failure		400		{object}	errResponse
//	@Failure		500		{object}	errResponse
//	@Failure		default	{object}	errResponse
//	@Router			/api/projects/{id} [delete]
func (h *Handler) DeleteById(ctx *gin.Context) {
	paramId := ctx.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		newErrResponse(ctx, http.StatusBadRequest, fmt.Sprintf("%s: message: invalid id param", err))
		return
	}

	if err = h.service.ProjectService.DeleteById(int64(id)); err != nil {
		newErrResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	h.cache.Delete("all")
	h.cache.Delete(paramId)

	ctx.JSON(http.StatusOK, statusResponse{"ok"})
}
