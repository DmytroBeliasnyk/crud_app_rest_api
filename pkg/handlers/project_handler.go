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
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		dto.ProjectDTO	true	"project info"
//	@Success		201		{integer}	integer			id
//	@Failure		400		{object}	errResponse
//	@Failure		500		{object}	errResponse
//	@Failure		default	{object}	errResponse
//	@Router			/api/projects/ [post]
func (h *Handler) create(c *gin.Context) {
	var input dto.ProjectDTO
	if err := c.BindJSON(&input); err != nil {
		newErrResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId := c.GetInt64("user_id")
	if userId == 0 {
		newErrResponse(c, http.StatusUnauthorized, "user not found")
		return
	}

	projectId, err := h.service.ProjectService.Create(input, userId)
	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.cache.Delete(fmt.Sprintf("all%d", userId))

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": projectId,
	})
}

// GetById godoc
//
//	@Summary		GetById
//	@Description	get project by id
//	@Tags			projects
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		query		integer	true	"project id"
//	@Success		200		{object}	dto.ProjectDTO
//	@Failure		400		{object}	errResponse
//	@Failure		500		{object}	errResponse
//	@Failure		default	{object}	errResponse
//	@Router			/api/projects [get]
func (h *Handler) getById(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Query("id"))
	if err != nil || projectId == 0 {
		newErrResponse(c, http.StatusBadRequest, fmt.Sprintf("%s: message: invalid id param", err))
		return
	}

	userId := c.GetInt64("user_id")
	if userId == 0 {
		newErrResponse(c, http.StatusUnauthorized, "user unauthorized")
		return
	}

	cache := fmt.Sprintf("%d%d", projectId, userId)

	project, err := h.cache.Get(cache)
	if err != nil {
		project, err = h.service.ProjectService.GetById(int64(projectId), userId)
		if err != nil {
			newErrResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		if err = h.cache.Set(cache, project, time.Hour); err != nil {
			h.cache.Delete(cache)
			newErrResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, project)
}

// GetAll godoc
//
//	@Summary		GetAll
//	@Description	get all projects
//	@Tags			projects
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Success		200		{array}		dto.ProjectDTO
//	@Failure		500		{object}	errResponse
//	@Failure		default	{object}	errResponse
//	@Router			/api/projects/ [get]
func (h *Handler) getAll(c *gin.Context) {
	userId := c.GetInt64("user_id")
	if userId == 0 {
		newErrResponse(c, http.StatusUnauthorized, "user unauthorized")
		return
	}

	cache := fmt.Sprintf("all%d", userId)

	projects, err := h.cache.Get(cache)
	if err != nil {
		projects, err = h.service.ProjectService.GetAll(userId)
		if err != nil {
			newErrResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		if err = h.cache.Set(cache, projects, time.Hour); err != nil {
			h.cache.Delete(cache)
			newErrResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, projects)
}

// UpdateById godoc
//
//	@Summary		UpdateById
//	@Description	update project by id
//	@Tags			projects
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		query		integer			true	"project id"
//	@Param			input	body		dto.ProjectDTO	true	"project info"
//	@Success		200		{object}	statusResponse
//	@Failure		400		{object}	errResponse
//	@Failure		500		{object}	errResponse
//	@Failure		default	{object}	errResponse
//	@Router			/api/projects [post]
func (h *Handler) updateById(c *gin.Context) {
	userId := c.GetInt64("user_id")
	if userId == 0 {
		newErrResponse(c, http.StatusUnauthorized, "user unauthorized")
		return
	}

	projectId, err := strconv.Atoi(c.Query("id"))
	if err != nil || projectId == 0 {
		newErrResponse(c, http.StatusBadRequest, fmt.Sprintf("%s: message: invalid id param", err))
		return
	}

	var input dto.UpdateProjectDTO
	if err := c.BindJSON(&input); err != nil {
		newErrResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		newErrResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.service.ProjectService.UpdateById(int64(projectId), input, userId); err != nil {
		newErrResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.cache.Delete(fmt.Sprintf("%d%d", projectId, userId))
	h.cache.Delete(fmt.Sprintf("all%d", userId))

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// DeleteById godoc
//
//	@Summary		DeleteById
//	@Description	delete project by id
//	@Tags			projects
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		query		integer	true	"project id"
//	@Success		200		{object}	statusResponse
//	@Failure		400		{object}	errResponse
//	@Failure		500		{object}	errResponse
//	@Failure		default	{object}	errResponse
//	@Router			/api/projects [delete]
func (h *Handler) deleteById(c *gin.Context) {
	userId := c.GetInt64("user_id")
	if userId == 0 {
		newErrResponse(c, http.StatusUnauthorized, "user unauthorized")
		return
	}

	projectId, err := strconv.Atoi(c.Query("id"))
	if err != nil || projectId == 0 {
		newErrResponse(c, http.StatusBadRequest, fmt.Sprintf("%s: message: invalid id param", err))
		return
	}

	if err = h.service.ProjectService.DeleteById(int64(projectId), userId); err != nil {
		newErrResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.cache.Delete(fmt.Sprintf("%d%d", projectId, userId))
	h.cache.Delete(fmt.Sprintf("all%d", userId))

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
