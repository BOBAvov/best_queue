package handler

import (
	"net/http"
	"sso/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getUserProfile возвращает профиль текущего пользователя
func (h *Handler) getUserProfile(c *gin.Context) {
	userId, ok := c.Get(userCtx)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user id not found in context"})
		return
	}

	user, err := h.service.GetUserByID(userId.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// getUsers возвращает список всех пользователей (только для админов)
func (h *Handler) getUsers(c *gin.Context) {
	isAdmin, ok := c.Get(userIsAdmin)
	if !ok || !isAdmin.(bool) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	users, err := h.service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// updateUser обновляет данные пользователя
func (h *Handler) updateUser(c *gin.Context) {
	userId, ok := c.Get(userCtx)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user id not found in context"})
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Пользователь может обновлять только свои данные, админ - любые
	targetUserID := userId.(int)
	isAdmin, ok := c.Get(userIsAdmin)
	if ok && isAdmin.(bool) {
		// Админ может указать ID пользователя для обновления
		if input.ID != 0 {
			targetUserID = input.ID
		}
	}

	err := h.service.UpdateUser(targetUserID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

// deleteUser удаляет пользователя (только для админов)
func (h *Handler) deleteUser(c *gin.Context) {
	isAdmin, ok := c.Get(userIsAdmin)
	if !ok || !isAdmin.(bool) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	err = h.service.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
