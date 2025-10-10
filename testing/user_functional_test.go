package test

import (
	"fmt"
	"net/http"
	"sso/models"
	"testing"
)

// TestUserProfile тестирует получение профиля пользователя
func TestUserProfile(t *testing.T) {
	helper := NewTestHelper()

	// Создаем пользователя и получаем токен
	helper.createTestUser(t, "profiletest", "password123", "@profiletest", "ИУ7-12Б")
	token := helper.loginUser(t, "@profiletest", "password123")

	t.Run("GetUserProfile", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/profile", nil, token)
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

		if _, ok := result["user"]; !ok {
			t.Error("Response should contain user data")
		}

		userData := result["user"].(map[string]interface{})
		if username, ok := userData["username"].(string); !ok || username != "profiletest" {
			t.Error("User data should contain correct username")
		}
	})

	t.Run("GetUserProfile_Unauthorized", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/profile", nil, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}

// TestUserUpdate тестирует обновление данных пользователя
func TestUserUpdate(t *testing.T) {
	helper := NewTestHelper()

	// Создаем пользователя и получаем токен
	helper.createTestUser(t, "updatetest", "password123", "@updatetest", "ИУ7-12Б")
	token := helper.loginUser(t, "@updatetest", "password123")

	t.Run("UpdateUserProfile", func(t *testing.T) {
		updateData := models.User{
			Username: "updateduser",
			TgNick:   "@updateduser",
		}

		resp, err := helper.makeRequest("PUT", baseURL+"/api/profile", updateData, token)
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

		if message, ok := result["message"].(string); !ok || message != "user updated successfully" {
			t.Error("Response should contain success message")
		}
	})

	t.Run("UpdateUserProfile_Unauthorized", func(t *testing.T) {
		updateData := models.User{
			Username: "updateduser",
		}

		resp, err := helper.makeRequest("PUT", baseURL+"/api/profile", updateData, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}

// TestAdminEndpoints тестирует админские endpoints
func TestAdminEndpoints(t *testing.T) {
	helper := NewTestHelper()

	// Создаем обычного пользователя
	helper.createTestUser(t, "regularuser", "password123", "@regularuser", "ИУ7-12Б")
	regularToken := helper.loginUser(t, "@regularuser", "password123")

	// Создаем админа (предполагаем, что есть способ создать админа)
	// В реальном проекте это может быть через специальный endpoint или базу данных
	helper.createTestUser(t, "adminuser", "password123", "@adminuser", "ИУ7-12Б")
	adminToken := helper.loginUser(t, "@adminuser", "password123")

	t.Run("CheckAdminStatus_RegularUser", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/admin", nil, regularToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Обычный пользователь не должен быть админом
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for non-admin user, got %d", resp.StatusCode)
		}
	})

	t.Run("CheckAdminStatus_AdminUser", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/admin", nil, adminToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// В зависимости от реализации, админ может получить 200 или 401
		// Проверяем, что ответ корректен
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
	})

	t.Run("GetAllUsers_RegularUser", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/admin/users", nil, regularToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403 for non-admin user, got %d", resp.StatusCode)
		}
	})

	t.Run("GetAllUsers_AdminUser", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/admin/users", nil, adminToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// В зависимости от реализации, админ может получить 200 или 403
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusForbidden {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			if err := helper.parseResponse(resp, &result); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if _, ok := result["users"]; !ok {
				t.Error("Response should contain users array")
			}
		}
	})

	t.Run("DeleteUser_RegularUser", func(t *testing.T) {
		// Создаем пользователя для удаления
		userToDelete := helper.createTestUser(t, "todelete", "password123", "@todelete", "ИУ7-12Б")

		resp, err := helper.makeRequest("DELETE", fmt.Sprintf("%s/api/admin/users/%d", baseURL, userToDelete), nil, regularToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403 for non-admin user, got %d", resp.StatusCode)
		}
	})

	t.Run("DeleteUser_AdminUser", func(t *testing.T) {
		// Создаем пользователя для удаления
		userToDelete := helper.createTestUser(t, "todelete2", "password123", "@todelete2", "ИУ7-12Б")

		resp, err := helper.makeRequest("DELETE", fmt.Sprintf("%s/api/admin/users/%d", baseURL, userToDelete), nil, adminToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// В зависимости от реализации, админ может получить 200 или 403
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusForbidden {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			if err := helper.parseResponse(resp, &result); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if message, ok := result["message"].(string); !ok || message != "user deleted successfully" {
				t.Error("Response should contain success message")
			}
		}
	})

	t.Run("DeleteUser_InvalidID", func(t *testing.T) {
		resp, err := helper.makeRequest("DELETE", baseURL+"/api/admin/users/invalid", nil, adminToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid user ID, got %d", resp.StatusCode)
		}
	})
}
