package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sso/models"
	"testing"
	"time"
)

const (
	baseURL = "http://127.0.0.1:8080"
)

// TestHelper содержит вспомогательные методы для тестов
type TestHelper struct {
	client *http.Client
}

func NewTestHelper() *TestHelper {
	return &TestHelper{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// makeRequest выполняет HTTP запрос
func (h *TestHelper) makeRequest(method, url string, body interface{}, token string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return h.client.Do(req)
}

// parseResponse парсит JSON ответ
func (h *TestHelper) parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

// createTestUser создает тестового пользователя
func (h *TestHelper) createTestUser(t *testing.T, username, password, tgNick, group string) int {
	userData := models.RegisterUser{
		Username: username,
		Password: password,
		TgNick:   tgNick,
		Group:    group,
	}

	resp, err := h.makeRequest("POST", baseURL+"/auth/sign-up", userData, "")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create test user, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := h.parseResponse(resp, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	id, ok := result["id"].(float64)
	if !ok {
		t.Fatalf("Invalid user ID in response")
	}

	return int(id)
}

// loginUser выполняет вход пользователя
func (h *TestHelper) loginUser(t *testing.T, tgNick, password string) string {
	authData := models.AuthUser{
		TgNick:   tgNick,
		Password: password,
	}

	resp, err := h.makeRequest("POST", baseURL+"/auth/sign-in", authData, "")
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to login user, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := h.parseResponse(resp, &result); err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}

	token, ok := result["token"].(string)
	if !ok {
		t.Fatalf("Invalid token in response")
	}

	return token
}

// createTestGroup создает тестовую группу (требует админ токен)
func (h *TestHelper) createTestGroup(t *testing.T, token, code, comment string) int {
	groupData := models.GroupCreateRequest{
		Code:    code,
		Comment: comment,
	}

	resp, err := h.makeRequest("POST", baseURL+"/api/groups", groupData, token)
	if err != nil {
		t.Fatalf("Failed to create test group: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create test group, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := h.parseResponse(resp, &result); err != nil {
		t.Fatalf("Failed to parse group response: %v", err)
	}

	id, ok := result["id"].(float64)
	if !ok {
		t.Fatalf("Invalid group ID in response")
	}

	return int(id)
}

// createTestQueue создает тестовую очередь (требует админ токен)
func (h *TestHelper) createTestQueue(t *testing.T, token, title string) int {
	queueData := models.CreateQueueRequest{
		Title:     title,
		TimeStart: time.Now(),
		TimeEnd:   time.Now().Add(2 * time.Hour),
	}

	resp, err := h.makeRequest("POST", baseURL+"/api/queues", queueData, token)
	if err != nil {
		t.Fatalf("Failed to create test queue: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create test queue, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := h.parseResponse(resp, &result); err != nil {
		t.Fatalf("Failed to parse queue response: %v", err)
	}

	id, ok := result["id"].(float64)
	if !ok {
		t.Fatalf("Invalid queue ID in response")
	}

	return int(id)
}
