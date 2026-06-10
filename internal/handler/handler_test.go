package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// Константы для контекста
type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

func TestRoomHandler_CreateRoom_InvalidBody(t *testing.T) {
	handler := NewRoomHandler(nil)

	req := httptest.NewRequest("POST", "/rooms/create", bytes.NewBufferString(`invalid json`))
	w := httptest.NewRecorder()

	handler.RoomCreate(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoomHandler_CreateRoom_EmptyName(t *testing.T) {
	handler := NewRoomHandler(nil)

	body := bytes.NewBufferString(`{"name": ""}`)
	req := httptest.NewRequest("POST", "/rooms/create", body)
	w := httptest.NewRecorder()

	handler.RoomCreate(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestScheduleHandler_CreateSchedule_InvalidBody(t *testing.T) {
	handler := NewScheduleHandler(nil)

	req := httptest.NewRequest("POST", "/rooms/123/schedule/create", bytes.NewBufferString(`invalid`))
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"roomId": "123"})

	handler.CreateSchedule(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestScheduleHandler_CreateSchedule_MissingFields(t *testing.T) {
	handler := NewScheduleHandler(nil)

	tests := []struct {
		name string
		body string
	}{
		{"missing daysOfWeek", `{"startTime":"09:00","endTime":"18:00"}`},
		{"missing startTime", `{"daysOfWeek":[1,2,3],"endTime":"18:00"}`},
		{"missing endTime", `{"daysOfWeek":[1,2,3],"startTime":"09:00"}`},
		{"empty daysOfWeek", `{"daysOfWeek":[],"startTime":"09:00","endTime":"18:00"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/rooms/123/schedule/create", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{"roomId": "123"})

			handler.CreateSchedule(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestSlotHandler_ListAvailableSlots_MissingDate(t *testing.T) {
	handler := NewSlotHandler(nil, nil)

	req := httptest.NewRequest("GET", "/rooms/123/slots/list", nil)
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"roomId": "123"})

	handler.ListAvailableSlots(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSlotHandler_ListAvailableSlots_InvalidDate(t *testing.T) {
	handler := NewSlotHandler(nil, nil)

	req := httptest.NewRequest("GET", "/rooms/123/slots/list?date=invalid-date", nil)
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"roomId": "123"})

	handler.ListAvailableSlots(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
