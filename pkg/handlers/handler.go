package handlers

import (
	"time"

	_ "github.com/DmytroBeliasnyk/crud_app_rest_api/docs"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/config"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock.go
type Cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string)
}

type Handler struct {
	service *services.AbstractService
	cfg     cookieConfig
	cache   Cache
}

type cookieConfig struct {
	name     string
	age      int
	path     string
	domain   string
	secure   bool
	httpOnly bool
}

func NewHandler(service *services.AbstractService, config *config.Config, cache Cache) *Handler {
	cooks := config.Cookie
	return &Handler{
		service: service,
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
			projects.POST("", h.updateById)
			projects.DELETE("", h.deleteById)
		}
	}

	return router
}
