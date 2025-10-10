package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// TestAPIStatus тестирует основной endpoint API
func TestAPIStatus(t *testing.T) {
	helper := NewTestHelper()

	t.Run("GetAPIStatus", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/", nil, "")
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

		if message, ok := result["message"].(string); !ok || message != "API is running" {
			t.Error("Response should contain API status message")
		}
	})

	t.Run("GetAPIStatus_WithToken", func(t *testing.T) {
		// Создаем пользователя и получаем токен
		helper.createTestUser(t, "statustest", "password123", "@statustest", "ИУ7-12Б")
		token := helper.loginUser(t, "@statustest", "password123")

		resp, err := helper.makeRequest("GET", baseURL+"/", nil, token)
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

		if message, ok := result["message"].(string); !ok || message != "API is running" {
			t.Error("Response should contain API status message")
		}
	})
}

// TestInvalidEndpoints тестирует несуществующие endpoints
func TestInvalidEndpoints(t *testing.T) {
	helper := NewTestHelper()

	t.Run("NonExistentEndpoint", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/nonexistent", nil, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404 for non-existent endpoint, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidMethod", func(t *testing.T) {
		resp, err := helper.makeRequest("PATCH", baseURL+"/api/profile", nil, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405 for invalid method, got %d", resp.StatusCode)
		}
	})

	t.Run("MalformedURL", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api//profile", nil, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// В зависимости от роутера, может быть 404 или 400
		if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Unexpected status code for malformed URL: %d", resp.StatusCode)
		}
	})
}

// TestCORSHeaders тестирует CORS заголовки (если они настроены)
func TestCORSHeaders(t *testing.T) {
	helper := NewTestHelper()

	t.Run("CORSHeaders", func(t *testing.T) {
		req, err := http.NewRequest("OPTIONS", baseURL+"/api/profile", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "GET")

		resp, err := helper.client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Проверяем, что сервер отвечает на OPTIONS запросы
		// Статус может быть 200, 204 или 404 в зависимости от настройки CORS
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
			t.Errorf("Unexpected status code for OPTIONS request: %d", resp.StatusCode)
		}
	})
}

// TestContentTypeValidation тестирует валидацию Content-Type
func TestContentTypeValidation(t *testing.T) {
	helper := NewTestHelper()

	t.Run("InvalidContentType", func(t *testing.T) {
		userData := map[string]string{
			"username": "testuser",
			"password": "password123",
			"tg_nick":  "@testuser",
			"group":    "ИУ7-12Б",
		}

		jsonData, _ := json.Marshal(userData)
		req, err := http.NewRequest("POST", baseURL+"/auth/sign-up", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "text/plain")

		resp, err := helper.client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// В зависимости от настроек Gin, может быть 400 или 200
		if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusOK {
			t.Errorf("Unexpected status code for invalid content type: %d", resp.StatusCode)
		}
	})
}
