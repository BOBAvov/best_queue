package test

import (
	"fmt"
	"net/http"
	"sso/models"
	"testing"
	"time"
)

// TestQueueCRUD тестирует CRUD операции с очередями
func TestQueueCRUD(t *testing.T) {
	helper := NewTestHelper()

	// Создаем админа и обычного пользователя
	helper.createTestUser(t, "queueadmin", "password123", "@queueadmin", "ИУ7-12Б")
	adminToken := helper.loginUser(t, "@queueadmin", "password123")

	helper.createTestUser(t, "queueuser", "password123", "@queueuser", "ИУ7-12Б")
	userToken := helper.loginUser(t, "@queueuser", "password123")

	t.Run("CreateQueue_Admin", func(t *testing.T) {
		queueData := models.CreateQueueRequest{
			Title:     "Test Queue",
			TimeStart: time.Now(),
			TimeEnd:   time.Now().Add(2 * time.Hour),
		}

		resp, err := helper.makeRequest("POST", baseURL+"/api/queues", queueData, adminToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := helper.parseResponse(resp, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if _, ok := result["id"]; !ok {
			t.Error("Response should contain queue ID")
		}

		if message, ok := result["message"].(string); !ok || message != "queue created successfully" {
			t.Error("Response should contain success message")
		}
	})

	t.Run("CreateQueue_RegularUser", func(t *testing.T) {
		queueData := models.CreateQueueRequest{
			Title:     "Test Queue 2",
			TimeStart: time.Now(),
			TimeEnd:   time.Now().Add(2 * time.Hour),
		}

		resp, err := helper.makeRequest("POST", baseURL+"/api/queues", queueData, userToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403 for non-admin user, got %d", resp.StatusCode)
		}
	})

	t.Run("CreateQueue_Unauthorized", func(t *testing.T) {
		queueData := models.CreateQueueRequest{
			Title:     "Test Queue 3",
			TimeStart: time.Now(),
			TimeEnd:   time.Now().Add(2 * time.Hour),
		}

		resp, err := helper.makeRequest("POST", baseURL+"/api/queues", queueData, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for unauthorized user, got %d", resp.StatusCode)
		}
	})

	t.Run("GetAllQueues", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/queues", nil, userToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := helper.parseResponse(resp, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if _, ok := result["queues"]; !ok {
			t.Error("Response should contain queues array")
		}
	})

	t.Run("GetQueueByID", func(t *testing.T) {
		// Создаем очередь для получения
		queueID := helper.createTestQueue(t, adminToken, "Get Queue Test")

		resp, err := helper.makeRequest("GET", fmt.Sprintf("%s/api/queues/%d", baseURL, queueID), nil, userToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := helper.parseResponse(resp, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if _, ok := result["queue"]; !ok {
			t.Error("Response should contain queue data")
		}
	})

	t.Run("GetQueueByID_NotFound", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/queues/99999", nil, userToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("GetQueueByID_InvalidID", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/queues/invalid", nil, userToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("UpdateQueue_Admin", func(t *testing.T) {
		// Создаем очередь для обновления
		queueID := helper.createTestQueue(t, adminToken, "Update Queue Test")

		updateData := models.Queue{
			ID:        queueID,
			Title:     "Updated Queue Title",
			TimeStart: time.Now(),
			TimeEnd:   time.Now().Add(3 * time.Hour),
		}

		resp, err := helper.makeRequest("PUT", fmt.Sprintf("%s/api/queues/%d", baseURL, queueID), updateData, adminToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := helper.parseResponse(resp, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if message, ok := result["message"].(string); !ok || message != "queue updated successfully" {
			t.Error("Response should contain success message")
		}
	})

	t.Run("UpdateQueue_RegularUser", func(t *testing.T) {
		// Создаем очередь для обновления
		queueID := helper.createTestQueue(t, adminToken, "Update Queue Test 2")

		updateData := models.Queue{
			ID:        queueID,
			Title:     "Updated Queue Title",
			TimeStart: time.Now(),
			TimeEnd:   time.Now().Add(3 * time.Hour),
		}

		resp, err := helper.makeRequest("PUT", fmt.Sprintf("%s/api/queues/%d", baseURL, queueID), updateData, userToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403 for non-admin user, got %d", resp.StatusCode)
		}
	})

	t.Run("DeleteQueue_Admin", func(t *testing.T) {
		// Создаем очередь для удаления
		queueID := helper.createTestQueue(t, adminToken, "Delete Queue Test")

		resp, err := helper.makeRequest("DELETE", fmt.Sprintf("%s/api/queues/%d", baseURL, queueID), nil, adminToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := helper.parseResponse(resp, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if message, ok := result["message"].(string); !ok || message != "queue deleted successfully" {
			t.Error("Response should contain success message")
		}
	})

	t.Run("DeleteQueue_RegularUser", func(t *testing.T) {
		// Создаем очередь для удаления
		queueID := helper.createTestQueue(t, adminToken, "Delete Queue Test 2")

		resp, err := helper.makeRequest("DELETE", fmt.Sprintf("%s/api/queues/%d", baseURL, queueID), nil, userToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403 for non-admin user, got %d", resp.StatusCode)
		}
	})
}

// TestQueueParticipation тестирует участие в очередях
func TestQueueParticipation(t *testing.T) {
	helper := NewTestHelper()

	// Создаем админа и пользователей
	helper.createTestUser(t, "queueadmin2", "password123", "@queueadmin2", "ИУ7-12Б")
	adminToken := helper.loginUser(t, "@queueadmin2", "password123")

	helper.createTestUser(t, "queueuser1", "password123", "@queueuser1", "ИУ7-12Б")
	user1Token := helper.loginUser(t, "@queueuser1", "password123")

	helper.createTestUser(t, "queueuser2", "password123", "@queueuser2", "ИУ7-12Б")
	user2Token := helper.loginUser(t, "@queueuser2", "password123")

	// Создаем очередь
	queueID := helper.createTestQueue(t, adminToken, "Participation Test Queue")

	t.Run("JoinQueue", func(t *testing.T) {
		joinData := models.JoinQueueRequest{
			QueueID: queueID,
		}

		resp, err := helper.makeRequest("POST", fmt.Sprintf("%s/api/queues/%d/join", baseURL, queueID), joinData, user1Token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := helper.parseResponse(resp, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if _, ok := result["participant_id"]; !ok {
			t.Error("Response should contain participant ID")
		}

		if message, ok := result["message"].(string); !ok || message != "joined queue successfully" {
			t.Error("Response should contain success message")
		}
	})

	t.Run("JoinQueue_AlreadyJoined", func(t *testing.T) {
		joinData := models.JoinQueueRequest{
			QueueID: queueID,
		}

		// Пытаемся присоединиться к очереди повторно
		resp, err := helper.makeRequest("POST", fmt.Sprintf("%s/api/queues/%d/join", baseURL, queueID), joinData, user1Token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Должна быть ошибка (400 или 409 в зависимости от реализации)
		if resp.StatusCode == http.StatusOK {
			t.Error("Expected error for already joined queue, got success")
		}
	})

	t.Run("JoinQueue_Unauthorized", func(t *testing.T) {
		joinData := models.JoinQueueRequest{
			QueueID: queueID,
		}

		resp, err := helper.makeRequest("POST", fmt.Sprintf("%s/api/queues/%d/join", baseURL, queueID), joinData, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("GetQueueParticipants", func(t *testing.T) {
		// Добавляем второго пользователя в очередь
		joinData := models.JoinQueueRequest{
			QueueID: queueID,
		}
		helper.makeRequest("POST", fmt.Sprintf("%s/api/queues/%d/join", baseURL, queueID), joinData, user2Token)

		resp, err := helper.makeRequest("GET", fmt.Sprintf("%s/api/queues/%d/participants", baseURL, queueID), nil, user1Token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := helper.parseResponse(resp, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if _, ok := result["participants"]; !ok {
			t.Error("Response should contain participants array")
		}
	})

	t.Run("LeaveQueue", func(t *testing.T) {
		resp, err := helper.makeRequest("DELETE", fmt.Sprintf("%s/api/queues/%d/leave", baseURL, queueID), nil, user1Token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := helper.parseResponse(resp, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if message, ok := result["message"].(string); !ok || message != "left queue successfully" {
			t.Error("Response should contain success message")
		}
	})

	t.Run("LeaveQueue_NotInQueue", func(t *testing.T) {
		// Пытаемся покинуть очередь, в которой не участвуем
		resp, err := helper.makeRequest("DELETE", fmt.Sprintf("%s/api/queues/%d/leave", baseURL, queueID), nil, user1Token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Должна быть ошибка (400 в зависимости от реализации)
		if resp.StatusCode == http.StatusOK {
			t.Error("Expected error for not being in queue, got success")
		}
	})

	t.Run("ShiftQueue_Admin", func(t *testing.T) {
		// Добавляем пользователя обратно в очередь
		joinData := models.JoinQueueRequest{
			QueueID: queueID,
		}
		helper.makeRequest("POST", fmt.Sprintf("%s/api/queues/%d/join", baseURL, queueID), joinData, user1Token)

		resp, err := helper.makeRequest("POST", fmt.Sprintf("%s/api/queues/%d/shift", baseURL, queueID), nil, adminToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := helper.parseResponse(resp, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if message, ok := result["message"].(string); !ok || message != "queue shifted successfully" {
			t.Error("Response should contain success message")
		}
	})

	t.Run("ShiftQueue_RegularUser", func(t *testing.T) {
		resp, err := helper.makeRequest("POST", fmt.Sprintf("%s/api/queues/%d/shift", baseURL, queueID), nil, user2Token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403 for non-admin user, got %d", resp.StatusCode)
		}
	})
}
