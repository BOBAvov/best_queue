// Package handler содержит HTTP обработчики для работы с очередями
package handler

import (
	"net/http"
	"sso/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// createQueue создает новую очередь (только для админов)
func (h *Handler) createQueue(c *gin.Context) {
	isAdmin, ok := c.Get(userIsAdmin)
	if !ok || !isAdmin.(bool) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	var input models.CreateQueueRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queue := models.Queue{
		Title:     input.Title,
		TimeStart: input.TimeStart,
		TimeEnd:   input.TimeEnd,
	}

	id, err := h.service.CreateQueue(queue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "message": "queue created successfully"})
}

// getQueue возвращает очередь по ID
func (h *Handler) getQueue(c *gin.Context) {
	idStr := c.Param("id")
	queueID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid queue id"})
		return
	}

	queue, err := h.service.GetQueueByID(queueID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "queue not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"queue": queue})
}

// getAllQueues возвращает все очереди
func (h *Handler) getAllQueues(c *gin.Context) {
	queues, err := h.service.GetAllQueues()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"queues": queues})
}

// updateQueue обновляет очередь (только для админов)
func (h *Handler) updateQueue(c *gin.Context) {
	isAdmin, ok := c.Get(userIsAdmin)
	if !ok || !isAdmin.(bool) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	id := c.Param("id")
	queueID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid queue id"})
		return
	}

	var input models.Queue
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = queueID
	err = h.service.UpdateQueue(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "queue updated successfully"})
}

// deleteQueue удаляет очередь (только для админов)
func (h *Handler) deleteQueue(c *gin.Context) {
	isAdmin, ok := c.Get(userIsAdmin)
	if !ok || !isAdmin.(bool) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	idStr := c.Param("id")
	queueID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid queue id"})
		return
	}

	err = h.service.DeleteQueue(queueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "queue deleted successfully"})
}

// joinQueue добавляет пользователя в очередь
func (h *Handler) joinQueue(c *gin.Context) {
	userId, ok := c.Get(userCtx)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user id not found in context"})
		return
	}

	var input models.JoinQueueRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	participantID, err := h.service.JoinQueue(input.QueueID, userId.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"participant_id": participantID, "message": "joined queue successfully"})
}

// leaveQueue удаляет пользователя из очереди
func (h *Handler) leaveQueue(c *gin.Context) {
	userId, ok := c.Get(userCtx)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user id not found in context"})
		return
	}

	queueIDStr := c.Param("id")
	queueID, err := strconv.Atoi(queueIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid queue id"})
		return
	}

	err = h.service.LeaveQueue(queueID, userId.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "left queue successfully"})
}

// getQueueParticipants возвращает участников очереди
func (h *Handler) getQueueParticipants(c *gin.Context) {
	queueIDStr := c.Param("id")
	queueID, err := strconv.Atoi(queueIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid queue id"})
		return
	}

	participants, err := h.service.GetQueueParticipants(queueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"participants": participants})
}

// shiftQueue сдвигает очередь (удаляет первого пользователя) - только для админов
func (h *Handler) shiftQueue(c *gin.Context) {
	isAdmin, ok := c.Get(userIsAdmin)
	if !ok || !isAdmin.(bool) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	queueIDStr := c.Param("id")
	queueID, err := strconv.Atoi(queueIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid queue id"})
		return
	}

	err = h.service.ShiftQueue(queueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "queue shifted successfully"})
}
