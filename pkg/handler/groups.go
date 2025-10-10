// Package handler содержит HTTP обработчики для работы с группами
package handler

import (
	"net/http"
	"sso/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// createGroup создает новую группу (только для админов)
func (h *Handler) createGroup(c *gin.Context) {
	var input models.GroupCreateRequest
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if isAdmin, exists := c.Get(userIsAdmin); !exists || !isAdmin.(bool) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := h.service.CreateGroup(input.Code, input.Comment)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "group created"})
}

// getAllGroups возвращает все группы
func (h *Handler) getAllGroups(c *gin.Context) {
	groups, err := h.service.GetAllGroups()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}

// getGroupByID возвращает группу по ID
func (h *Handler) getGroupByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	group, err := h.service.GetGroupByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, group)
}

// updateGroup обновляет информацию о группе (только для админов)
func (h *Handler) updateGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	var input models.Group
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if isAdmin, exists := c.Get(userIsAdmin); !exists || !isAdmin.(bool) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.service.UpdateGroup(id, input); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "group updated"})
}

// deleteGroup удаляет группу (только для админов)
func (h *Handler) deleteGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	if isAdmin, exists := c.Get(userIsAdmin); !exists || !isAdmin.(bool) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.service.DeleteGroup(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "group deleted"})
}
