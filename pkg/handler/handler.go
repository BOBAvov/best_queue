// Package handler содержит HTTP обработчики для REST API
package handler

import (
	"net/http"
	"sso/models"
	"sso/pkg/services"
	"strings"

	"github.com/gin-gonic/gin"
)

// Константы для работы с контекстом и заголовками
const (
	authorizationHeader = "Authorization" // Название заголовка авторизации
	userCtx             = "userId"        // Ключ для хранения ID пользователя в контексте
	userIsAdmin         = "isAdmin"       // Ключ для хранения флага администратора в контексте
)

// Handler содержит сервисы для обработки HTTP запросов
type Handler struct {
	service services.Authorization // Сервис авторизации и работы с пользователями
}

// NewHandler создает новый экземпляр обработчика с указанным сервисом авторизации
func NewHandler(authService services.Authorization) *Handler {
	return &Handler{service: authService}
}

// InitRoutes инициализирует и настраивает все маршруты API
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	// Основной endpoint для проверки статуса API
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is running"})
	})

	// Группа маршрутов для аутентификации (не требует авторизации)
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp) // Регистрация нового пользователя
		auth.POST("/sign-in", h.signIn) // Вход в систему
	}

	// Группа защищенных маршрутов (требует JWT токен)
	api := router.Group("/api", h.userIdentity)
	{
		// Маршруты для работы с пользователями
		api.GET("/admin", h.isAdmin)          // Проверка статуса администратора
		api.GET("/profile", h.getUserProfile) // Получение профиля пользователя
		api.PUT("/profile", h.updateUser)     // Обновление профиля пользователя

		// Маршруты только для администраторов
		admin := api.Group("/admin")
		{
			admin.GET("/users", h.getUsers)          // Получение списка всех пользователей
			admin.DELETE("/users/:id", h.deleteUser) // Удаление пользователя
		}

		// Маршруты для работы с очередями
		queues := api.Group("/queues")
		{
			queues.GET("/", h.getAllQueues)                         // Получение всех очередей
			queues.POST("/", h.createQueue)                         // Создание очереди (только админ)
			queues.GET("/:id", h.getQueue)                          // Получение очереди по ID
			queues.PUT("/:id", h.updateQueue)                       // Обновление очереди (только админ)
			queues.DELETE("/:id", h.deleteQueue)                    // Удаление очереди (только админ)
			queues.POST("/:id/join", h.joinQueue)                   // Присоединение к очереди
			queues.DELETE("/:id/leave", h.leaveQueue)               // Покидание очереди
			queues.GET("/:id/participants", h.getQueueParticipants) // Получение участников очереди
			queues.POST("/:id/shift", h.shiftQueue)                 // Сдвиг очереди (только админ)
		}

		// Маршруты для работы с группами
		groups := api.Group("/groups")
		{
			groups.POST("/", h.createGroup)      // Создание группы (только админ)
			groups.GET("/", h.getAllGroups)      // Получение всех групп
			groups.GET("/:id", h.getGroupByID)   // Получение группы по ID
			groups.PUT("/:id", h.updateGroup)    // Обновление группы (только админ)
			groups.DELETE("/:id", h.deleteGroup) // Удаление группы (только админ)
		}
	}

	return router
}

// isAdmin проверяет, является ли текущий пользователь администратором
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

// signUp обрабатывает регистрацию нового пользователя
func (h *Handler) signUp(c *gin.Context) {
	var input models.RegisterUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, что все обязательные поля заполнены
	if input.Username == "" || input.Password == "" || input.TgNick == "" || input.Group == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all fields are required"})
		return
	}

	// Создаем пользователя через сервис
	id, err := h.service.CreateUser(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "message": "successfully signed up"})
}

// signIn обрабатывает вход пользователя в систему
func (h *Handler) signIn(c *gin.Context) {
	var input models.AuthUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Генерируем JWT токен для пользователя
	token, err := h.service.NewToken(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// userIdentity middleware для проверки JWT токена и извлечения информации о пользователе
func (h *Handler) userIdentity(c *gin.Context) {
	// Получаем заголовок авторизации
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
		return
	}

	// Проверяем формат заголовка "Bearer <token>"
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
		return
	}

	// Проверяем, что токен не пустой
	if len(headerParts[1]) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
		return
	}

	// Парсим токен и извлекаем информацию о пользователе
	userId, isAdmin, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Сохраняем ID пользователя и статус администратора в контексте для дальнейшего использования
	c.Set(userCtx, userId)
	c.Set(userIsAdmin, isAdmin)
	c.Next()
}
