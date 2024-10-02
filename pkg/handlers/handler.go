package handlers

import (
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.AbstractService
}

func NewHandler(service *services.AbstractService) *Handler {
	return &Handler{service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		projects := api.Group("/projects")
		{
			projects.Handle("POST", "/create", h.Add)
			projects.Handle("GET", "/", h.GetAll)
			projects.Handle("GET", "/:id", h.GetById)
			projects.Handle("POST", "/:id", h.UpdateById)
			projects.Handle("DELETE", "/:id", h.DeleteById)
		}
	}

	return router
}
