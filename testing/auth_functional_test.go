package test

import (
	"net/http"
	"sso/models"
	"testing"
)

// TestAuthSignUp тестирует регистрацию пользователя
func TestAuthSignUp(t *testing.T) {
	helper := NewTestHelper()

	t.Run("ValidSignUp", func(t *testing.T) {
		userData := models.RegisterUser{
			Username: "testuser1",
			Password: "password123",
			TgNick:   "@testuser1",
			Group:    "ИУ7-12Б",
		}

		resp, err := helper.makeRequest("POST", baseURL+"/auth/sign-up", userData, "")
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
			t.Error("Response should contain user ID")
		}

		if message, ok := result["message"].(string); !ok || message != "successfully signed up" {
			t.Error("Response should contain success message")
		}
	})

	t.Run("InvalidSignUp_MissingFields", func(t *testing.T) {
		userData := models.RegisterUser{
			Username: "testuser2",
			// Missing password, tg_nick, group
		}

		resp, err := helper.makeRequest("POST", baseURL+"/auth/sign-up", userData, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidSignUp_EmptyFields", func(t *testing.T) {
		userData := models.RegisterUser{
			Username: "",
			Password: "",
			TgNick:   "",
			Group:    "",
		}

		resp, err := helper.makeRequest("POST", baseURL+"/auth/sign-up", userData, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidSignUp_DuplicateUser", func(t *testing.T) {
		// Сначала создаем пользователя
		userData := models.RegisterUser{
			Username: "duplicateuser",
			Password: "password123",
			TgNick:   "@duplicateuser",
			Group:    "ИУ7-12Б",
		}

		resp, err := helper.makeRequest("POST", baseURL+"/auth/sign-up", userData, "")
		if err != nil {
			t.Fatalf("Failed to make first request: %v", err)
		}
		resp.Body.Close()

		// Пытаемся создать того же пользователя снова
		resp, err = helper.makeRequest("POST", baseURL+"/auth/sign-up", userData, "")
		if err != nil {
			t.Fatalf("Failed to make second request: %v", err)
		}
		defer resp.Body.Close()

		// Должна быть ошибка (возможно 500 или 400 в зависимости от реализации)
		if resp.StatusCode == http.StatusOK {
			t.Error("Expected error for duplicate user, got success")
		}
	})
}

// TestAuthSignIn тестирует вход пользователя
func TestAuthSignIn(t *testing.T) {
	helper := NewTestHelper()

	// Создаем тестового пользователя
	userID := helper.createTestUser(t, "signintest", "password123", "@signintest", "ИУ7-12Б")

	t.Run("ValidSignIn", func(t *testing.T) {
		authData := models.AuthUser{
			TgNick:   "@signintest",
			Password: "password123",
		}

		resp, err := helper.makeRequest("POST", baseURL+"/auth/sign-in", authData, "")
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

		if _, ok := result["token"].(string); !ok {
			t.Error("Response should contain JWT token")
		}
	})

	t.Run("InvalidSignIn_WrongPassword", func(t *testing.T) {
		authData := models.AuthUser{
			TgNick:   "@signintest",
			Password: "wrongpassword",
		}

		resp, err := helper.makeRequest("POST", baseURL+"/auth/sign-in", authData, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidSignIn_NonExistentUser", func(t *testing.T) {
		authData := models.AuthUser{
			TgNick:   "@nonexistent",
			Password: "password123",
		}

		resp, err := helper.makeRequest("POST", baseURL+"/auth/sign-in", authData, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidSignIn_MissingFields", func(t *testing.T) {
		authData := models.AuthUser{
			TgNick: "@signintest",
			// Missing password
		}

		resp, err := helper.makeRequest("POST", baseURL+"/auth/sign-in", authData, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	// Очистка
	_ = userID // Используем переменную чтобы избежать предупреждения
}

// TestAuthTokenValidation тестирует валидацию токенов
func TestAuthTokenValidation(t *testing.T) {
	helper := NewTestHelper()

	// Создаем пользователя и получаем токен
	helper.createTestUser(t, "tokentest", "password123", "@tokentest", "ИУ7-12Б")
	token := helper.loginUser(t, "@tokentest", "password123")

	t.Run("ValidToken", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/profile", nil, token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidToken_Empty", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/profile", nil, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidToken_Malformed", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/profile", nil, "invalid_token")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidToken_WrongFormat", func(t *testing.T) {
		// Устанавливаем неправильный заголовок Authorization
		req, err := http.NewRequest("GET", baseURL+"/api/profile", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", "InvalidFormat "+token)

		resp, err := helper.client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}
