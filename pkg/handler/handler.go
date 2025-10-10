package handler

import (
	"fmt"
	"net/http"
	"sso/models"
	"sso/pkg/services"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	userIsAdmin         = "isAdmin"
)

type Handler struct {
	service services.Authorization
}

func NewHandler(authService services.Authorization) *Handler {
	return &Handler{service: authService}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is running"})
	})

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		// User routes
		api.GET("/admin", h.isAdmin)
		api.GET("/profile", h.getUserProfile)
		api.PUT("/profile", h.updateUser)

		// Admin only routes
		admin := api.Group("/admin")
		{
			admin.GET("/users", h.getUsers)
			admin.DELETE("/users/:id", h.deleteUser)
		}

		// Queue routes
		queues := api.Group("/queues")
		{
			queues.GET("/", h.getAllQueues)
			queues.POST("/", h.createQueue) // Admin only
			queues.GET("/:id", h.getQueue)
			queues.PUT("/:id", h.updateQueue)    // Admin only
			queues.DELETE("/:id", h.deleteQueue) // Admin only
			queues.POST("/:id/join", h.joinQueue)
			queues.DELETE("/:id/leave", h.leaveQueue)
			queues.GET("/:id/participants", h.getQueueParticipants)
			queues.POST("/:id/shift", h.shiftQueue) // Admin only
		}
	}

	return router
}
func (h *Handler) isAdmin(c *gin.Context) {
	isAdmin, ok := c.Get(userIsAdmin)
	if !ok || !isAdmin.(bool) {
		c.Set(userIsAdmin, false)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	isAdmin = isAdmin.(bool)
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"admin": isAdmin})
}

func (h *Handler) signUp(c *gin.Context) {
	var input models.RegisterUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Username == "" || input.Password == "" || input.TgNick == "" || input.Group == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all fields are required"})
		return
	}

	id, err := h.service.CreateUser(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "message": "successfully signed up"})
}

func (h *Handler) signIn(c *gin.Context) {
	var input models.AuthUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(input)
	token, err := h.service.NewToken(input)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
		return
	}

	if len(headerParts[1]) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
		return
	}

	userId, isAdmin, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Сохраняем ID пользователя в контексте для дальнейшего использования
	c.Set(userCtx, userId)
	c.Set(userIsAdmin, isAdmin)
	c.Next()
}
