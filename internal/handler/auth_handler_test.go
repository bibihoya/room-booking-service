package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_DummyLogin(t *testing.T) {
	handler := NewAuthHandler("test-secret", 24, nil)

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"admin role", `{"role":"admin"}`, http.StatusOK},
		{"user role", `{"role":"user"}`, http.StatusOK},
		{"invalid role", `{"role":"invalid"}`, http.StatusBadRequest},
		{"empty body", `{}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/dummyLogin", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			handler.DummyLogin(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAuthHandler_DummyLogin_WithContext(t *testing.T) {
	handler := NewAuthHandler("test-secret", 24, nil)

	req := httptest.NewRequest("POST", "/dummyLogin", bytes.NewBufferString(`{"role":"admin"}`))
	w := httptest.NewRecorder()
	handler.DummyLogin(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	token := response["token"].(string)
	assert.NotEmpty(t, token)
}
