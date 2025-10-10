package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"sso/models"
	"testing"
)

// TestGroupCRUD тестирует CRUD операции с группами
func TestGroupCRUD(t *testing.T) {
	helper := NewTestHelper()

	// Создаем админа и обычного пользователя
	helper.createTestUser(t, "groupadmin", "password123", "@groupadmin", "ИУ7-12Б")
	adminToken := helper.loginUser(t, "@groupadmin", "password123")

	helper.createTestUser(t, "groupuser", "password123", "@groupuser", "ИУ7-12Б")
	userToken := helper.loginUser(t, "@groupuser", "password123")

	t.Run("CreateGroup_Admin", func(t *testing.T) {
		groupData := models.GroupCreateRequest{
			Code:    "ИУ7-12Б",
			Comment: "Test group for functional testing",
		}

		resp, err := helper.makeRequest("POST", baseURL+"/api/groups", groupData, adminToken)
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
			t.Error("Response should contain group ID")
		}

		if message, ok := result["message"].(string); !ok || message != "group created" {
			t.Error("Response should contain success message")
		}
	})

	t.Run("CreateGroup_RegularUser", func(t *testing.T) {
		groupData := models.GroupCreateRequest{
			Code:    "ИУ7-13Б",
			Comment: "Test group for regular user",
		}

		resp, err := helper.makeRequest("POST", baseURL+"/api/groups", groupData, userToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for non-admin user, got %d", resp.StatusCode)
		}
	})

	t.Run("CreateGroup_Unauthorized", func(t *testing.T) {
		groupData := models.GroupCreateRequest{
			Code:    "ИУ7-14Б",
			Comment: "Test group for unauthorized user",
		}

		resp, err := helper.makeRequest("POST", baseURL+"/api/groups", groupData, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for unauthorized user, got %d", resp.StatusCode)
		}
	})

	t.Run("CreateGroup_InvalidData", func(t *testing.T) {
		groupData := models.GroupCreateRequest{
			Code: "", // Пустой код
		}

		resp, err := helper.makeRequest("POST", baseURL+"/api/groups", groupData, adminToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid data, got %d", resp.StatusCode)
		}
	})

	t.Run("GetAllGroups", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/groups", nil, userToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result []models.Group
		if err := helper.parseResponse(resp, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		// Проверяем, что получили массив групп
		if len(result) == 0 {
			t.Error("Should return at least one group")
		}
	})

	t.Run("GetAllGroups_Unauthorized", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/groups", nil, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for unauthorized user, got %d", resp.StatusCode)
		}
	})

	t.Run("GetGroupByID", func(t *testing.T) {
		// Создаем группу для получения
		groupID := helper.createTestGroup(t, adminToken, "ИУ7-15Б", "Test group for get by ID")

		resp, err := helper.makeRequest("GET", fmt.Sprintf("%s/api/groups/%d", baseURL, groupID), nil, userToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result models.Group
		if err := helper.parseResponse(resp, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if result.ID != groupID {
			t.Errorf("Expected group ID %d, got %d", groupID, result.ID)
		}

		if result.Code != "ИУ7-15Б" {
			t.Errorf("Expected group code ИУ7-15Б, got %s", result.Code)
		}
	})

	t.Run("GetGroupByID_NotFound", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/groups/99999", nil, userToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status 500 for non-existent group, got %d", resp.StatusCode)
		}
	})

	t.Run("GetGroupByID_InvalidID", func(t *testing.T) {
		resp, err := helper.makeRequest("GET", baseURL+"/api/groups/invalid", nil, userToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid group ID, got %d", resp.StatusCode)
		}
	})

	t.Run("UpdateGroup_Admin", func(t *testing.T) {
		// Создаем группу для обновления
		groupID := helper.createTestGroup(t, adminToken, "ИУ7-16Б", "Original comment")

		updateData := models.Group{
			ID:      groupID,
			Code:    "ИУ7-16Б-Updated",
			Comment: sql.NullString{String: "Updated comment", Valid: true},
		}

		resp, err := helper.makeRequest("PUT", fmt.Sprintf("%s/api/groups/%d", baseURL, groupID), updateData, adminToken)
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

		if message, ok := result["message"].(string); !ok || message != "group updated" {
			t.Error("Response should contain success message")
		}
	})

	t.Run("UpdateGroup_RegularUser", func(t *testing.T) {
		// Создаем группу для обновления
		groupID := helper.createTestGroup(t, adminToken, "ИУ7-17Б", "Original comment")

		updateData := models.Group{
			ID:      groupID,
			Code:    "ИУ7-17Б-Updated",
			Comment: sql.NullString{String: "Updated comment", Valid: true},
		}

		resp, err := helper.makeRequest("PUT", fmt.Sprintf("%s/api/groups/%d", baseURL, groupID), updateData, userToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for non-admin user, got %d", resp.StatusCode)
		}
	})

	t.Run("UpdateGroup_InvalidData", func(t *testing.T) {
		// Создаем группу для обновления
		groupID := helper.createTestGroup(t, adminToken, "ИУ7-18Б", "Original comment")

		updateData := models.Group{
			ID:   groupID,
			Code: "", // Пустой код
		}

		resp, err := helper.makeRequest("PUT", fmt.Sprintf("%s/api/groups/%d", baseURL, groupID), updateData, adminToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid data, got %d", resp.StatusCode)
		}
	})

	t.Run("DeleteGroup_Admin", func(t *testing.T) {
		// Создаем группу для удаления
		groupID := helper.createTestGroup(t, adminToken, "ИУ7-19Б", "Group to be deleted")

		resp, err := helper.makeRequest("DELETE", fmt.Sprintf("%s/api/groups/%d", baseURL, groupID), nil, adminToken)
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

		if message, ok := result["message"].(string); !ok || message != "group deleted" {
			t.Error("Response should contain success message")
		}
	})

	t.Run("DeleteGroup_RegularUser", func(t *testing.T) {
		// Создаем группу для удаления
		groupID := helper.createTestGroup(t, adminToken, "ИУ7-20Б", "Group to be deleted by regular user")

		resp, err := helper.makeRequest("DELETE", fmt.Sprintf("%s/api/groups/%d", baseURL, groupID), nil, userToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for non-admin user, got %d", resp.StatusCode)
		}
	})

	t.Run("DeleteGroup_NotFound", func(t *testing.T) {
		resp, err := helper.makeRequest("DELETE", baseURL+"/api/groups/99999", nil, adminToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status 500 for non-existent group, got %d", resp.StatusCode)
		}
	})

	t.Run("DeleteGroup_InvalidID", func(t *testing.T) {
		resp, err := helper.makeRequest("DELETE", baseURL+"/api/groups/invalid", nil, adminToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid group ID, got %d", resp.StatusCode)
		}
	})
}

// TestGroupEdgeCases тестирует граничные случаи для групп
func TestGroupEdgeCases(t *testing.T) {
	helper := NewTestHelper()

	// Создаем админа
	helper.createTestUser(t, "groupadmin2", "password123", "@groupadmin2", "ИУ7-12Б")
	adminToken := helper.loginUser(t, "@groupadmin2", "password123")

	t.Run("CreateGroup_DuplicateCode", func(t *testing.T) {
		groupData := models.GroupCreateRequest{
			Code:    "ИУ7-21Б",
			Comment: "First group with this code",
		}

		// Создаем первую группу
		resp, err := helper.makeRequest("POST", baseURL+"/api/groups", groupData, adminToken)
		if err != nil {
			t.Fatalf("Failed to make first request: %v", err)
		}
		resp.Body.Close()

		// Пытаемся создать группу с тем же кодом
		resp, err = helper.makeRequest("POST", baseURL+"/api/groups", groupData, adminToken)
		if err != nil {
			t.Fatalf("Failed to make second request: %v", err)
		}
		defer resp.Body.Close()

		// Должна быть ошибка (500 или 400 в зависимости от реализации)
		if resp.StatusCode == http.StatusOK {
			t.Error("Expected error for duplicate group code, got success")
		}
	})

	t.Run("CreateGroup_LongCode", func(t *testing.T) {
		longCode := "ОченьДлинныйКодГруппыКоторыйМожетПревышатьОбычныеОграничения"
		groupData := models.GroupCreateRequest{
			Code:    longCode,
			Comment: "Group with very long code",
		}

		resp, err := helper.makeRequest("POST", baseURL+"/api/groups", groupData, adminToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// В зависимости от ограничений БД, может быть 200 или 400
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
	})

	t.Run("CreateGroup_SpecialCharacters", func(t *testing.T) {
		groupData := models.GroupCreateRequest{
			Code:    "ИУ7-22Б@#$%",
			Comment: "Group with special characters",
		}

		resp, err := helper.makeRequest("POST", baseURL+"/api/groups", groupData, adminToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// В зависимости от валидации, может быть 200 или 400
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
	})
}
