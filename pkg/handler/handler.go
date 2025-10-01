package handler

import (
	"net/http"
	"sso/models"
	"sso/pkg/services"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

type Handler struct {
	service services.Authorization
}

func NewHandler(authService services.Authorization) *Handler {
	return &Handler{service: authService}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is running"})
	})

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity, h.isAdmin)
	{
		api.GET("/profile", h.getProfile)
	}

	create := router.Group("/create")
	{
		create.POST("/department")
		create.POST("/faculty")
		create.POST("/group", h.createGroup)
		create.POST("/admin", h.createAdmin)
	}

	return router
}

func (h *Handler) signUp(c *gin.Context) {
	var input models.RegisterUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Username == "" || input.Password == "" || input.TgNick == "" || input.Group == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "<UNK>"})
	}

	id, err := h.service.CreateUser(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "message": "successfully signed up"})
}

type signInInput struct {
	Tg_name  string `json:"tg_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//token, err := h.service.GenerateToken(input.Tg_name, input.Password)
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	//	return
	//}

	//c.JSON(http.StatusOK, gin.H{"token": token})
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

	userId, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Сохраняем ID пользователя в контексте для дальнейшего использования
	c.Set(userCtx, userId)
}

func (h *Handler) getProfile(c *gin.Context) {
	// Получаем ID пользователя из контекста, который был установлен в middleware
	userId, ok := c.Get(userCtx)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user id not found in context"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "this is a protected route", "user_id": userId})
}

func (h *Handler) isAdmin(c *gin.Context) {
	userId, ok := c.Get(userCtx)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user id not found in context"})
		return
	}

	c.Set(userCtx, userId)
}

// TODO: потом уберу, ручка для настройки
func (h *Handler) createAdmin(c *gin.Context) {
	var input models.RegisterUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//h.service.CreateGroup()
}

type createGroupInput struct {
	Name string `json:"name"`
}

// TODO: потом уберу, ручка для настройки
func (h *Handler) createGroup(c *gin.Context) {
	var input createGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.service.CreateGroup(input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "successfully created group"})
}

func (h *Handler) createFaculty(c *gin.Context) {
	var input models.Faculty
}
