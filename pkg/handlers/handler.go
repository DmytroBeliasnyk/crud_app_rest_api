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
	auth    *services.AuthService
	cache   *memory.Cache
}

func NewHandler(service *services.AbstractService, auth *services.AuthService, cache *memory.Cache) *Handler {
	return &Handler{
		service: service,
		auth:    auth,
		cache:   cache,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api")
	{
		api.Use(h.middlewareAuth)

		projects := api.Group("/projects")
		{
			projects.POST("/", h.create)
			projects.GET("/", h.getAll)
			projects.GET("", h.getById)
			projects.POST("/:id", h.updateById)
			projects.DELETE("/:id", h.deleteById)
		}
	}

	return router
}
