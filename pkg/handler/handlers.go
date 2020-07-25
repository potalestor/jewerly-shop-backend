package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/service"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h Handler) Init() *gin.Engine {
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAuthRoutes(router)
	h.initAPIRoutes(router)
	h.initAdminRoutes(router)

	return router
}

func (h *Handler) initAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
	}
}

func (h *Handler) initAPIRoutes(router *gin.Engine) {
	api := router.Group("/api", h.userIdentity)
	{
		user := api.Group("/user")
		{
			user.GET("/profile")
		}
	}
}

func (h *Handler) initAdminRoutes(router *gin.Engine) {
	_ = router.Group("/admin")
	{

	}
}