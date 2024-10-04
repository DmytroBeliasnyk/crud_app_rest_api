package handlers

import (
	_ "github.com/DmytroBeliasnyk/crud_app_rest_api/docs"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services"
	"github.com/DmytroBeliasnyk/in_memory_cache/memory"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *services.AbstractService
	cache   *memory.Cache
}

func NewHandler(service *services.AbstractService, cache *memory.Cache) *Handler {
	return &Handler{
		service: service,
		cache:   cache,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		projects := api.Group("/projects")
		{
			projects.POST("/", h.Create)
			projects.GET("/", h.GetAll)
			projects.GET("", h.GetById)
			projects.POST("/:id", h.UpdateById)
			projects.DELETE("/:id", h.DeleteById)
		}
	}

	return router
}
