package handlers

import (
	_ "github.com/DmytroBeliasnyk/crud_app_rest_api/docs"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/config"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services/implserv"
	"github.com/DmytroBeliasnyk/in_memory_cache/memory"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *services.AbstractService
	auth    *implserv.AuthService
	cfg     cookieConfig
	cache   *memory.Cache
}

type cookieConfig struct {
	name     string
	age      int
	path     string
	domain   string
	secure   bool
	httpOnly bool
}

func NewHandler(service *services.AbstractService, auth *implserv.AuthService,
	config *config.Config, cache *memory.Cache) *Handler {
	cooks := config.Cookie
	return &Handler{
		service: service,
		auth:    auth,
		cache:   cache,
		cfg: cookieConfig{
			name:     cooks.Name,
			age:      cooks.Age,
			path:     cooks.Path,
			domain:   cooks.Domain,
			secure:   cooks.Secure,
			httpOnly: cooks.HttpOnly,
		},
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.GET("/refresh", h.refresh)
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
