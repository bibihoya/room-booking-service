package domain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendError(t *testing.T) {
	w := httptest.NewRecorder()
	SendError(w, "TEST_ERROR", "test message", http.StatusBadRequest)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	errorObj := response["error"].(map[string]interface{})
	assert.Equal(t, "TEST_ERROR", errorObj["code"])
	assert.Equal(t, "test message", errorObj["message"])
}

func TestSendSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"key": "value"}
	SendSuccess(w, data, http.StatusCreated)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "value", response["key"])
}

func TestSendSuccessWithDifferentData(t *testing.T) {
	tests := []struct {
		name       string
		data       interface{}
		statusCode int
	}{
		{"string data", map[string]string{"result": "ok"}, http.StatusOK},
		{"array data", map[string][]int{"ids": {1, 2, 3}}, http.StatusOK},
		{"empty data", map[string]interface{}{}, http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			SendSuccess(w, tt.data, tt.statusCode)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
