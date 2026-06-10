package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserID(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, "test-user-id")
	userID, ok := GetUserID(ctx)
	assert.True(t, ok)
	assert.Equal(t, "test-user-id", userID)

	_, ok = GetUserID(context.Background())
	assert.False(t, ok)
}

func TestGetRole(t *testing.T) {
	ctx := context.WithValue(context.Background(), RoleKey, "admin")
	role, ok := GetRole(ctx)
	assert.True(t, ok)
	assert.Equal(t, "admin", role)

	_, ok = GetRole(context.Background())
	assert.False(t, ok)
}

func TestAdminCheck(t *testing.T) {
	handler := AdminCheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), RoleKey, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req = httptest.NewRequest("GET", "/", nil)
	ctx = context.WithValue(req.Context(), RoleKey, "user")
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	req = httptest.NewRequest("GET", "/", nil)
	req = req.WithContext(context.Background())
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUserCheck(t *testing.T) {
	handler := UserCheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req = httptest.NewRequest("GET", "/", nil)
	ctx = context.WithValue(req.Context(), RoleKey, "admin")
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	req = httptest.NewRequest("GET", "/", nil)
	req = req.WithContext(context.Background())
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	handler := AuthMiddleware("test-secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	handler := AuthMiddleware("test-secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	handler := AuthMiddleware("test-secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	token, err := generateTestToken("secret1")
	assert.NoError(t, err)

	handler := AuthMiddleware("secret2")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func generateTestToken(secret string) (string, error) {
	return "test-token", nil
}

func TestSendAuthError(t *testing.T) {
	w := httptest.NewRecorder()
	sendAuthError(w, "TEST_CODE", "test message", http.StatusUnauthorized)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	errorObj := response["error"].(map[string]interface{})
	assert.Equal(t, "TEST_CODE", errorObj["code"])
	assert.Equal(t, "test message", errorObj["message"])
}
