package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/domain"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/storage"
)

type ScheduleHandler struct {
	scheduleStorage *storage.ScheduleStorage
}

func NewScheduleHandler(scheduleStorage *storage.ScheduleStorage) *ScheduleHandler {
	return &ScheduleHandler{scheduleStorage}
}

// CreateSchedule godoc
// @Summary      Создать расписание переговорки
// @Description  Создаёт расписание для переговорки. Только admin, только один раз. Длительность слота 30 минут.
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        roomId path string true "ID переговорки"
// @Param        request body object true "Расписание" example({"daysOfWeek": [1,2,3,4,5], "startTime": "09:00", "endTime": "18:00"})
// @Success      201  {object}  map[string]interface{}  "schedule"
// @Failure      400  {object}  map[string]interface{}  "error"
// @Failure      401  {object}  map[string]interface{}  "error"
// @Failure      403  {object}  map[string]interface{}  "error"
// @Failure      404  {object}  map[string]interface{}  "error"
// @Failure      409  {object}  map[string]interface{}  "error"
// @Failure      500  {object}  map[string]interface{}  "error"
// @Router       /rooms/{roomId}/schedule/create [post]
func (h *ScheduleHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomId"]

	var req struct {
		DaysOfWeek []int  `json:"daysOfWeek"`
		StartTime  string `json:"startTime"`
		EndTime    string `json:"endTime"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.SendError(w, "INVALID_REQUEST", "Invalid request", http.StatusBadRequest)
		return
	}

	if len(req.DaysOfWeek) == 0 {
		domain.SendError(w, "INVALID_REQUEST", "daysOfWeek is required", http.StatusBadRequest)
		return
	}

	if req.StartTime == "" {
		domain.SendError(w, "INVALID_REQUEST", "start time is required", http.StatusBadRequest)
		return
	}

	if req.EndTime == "" {
		domain.SendError(w, "INVALID_REQUEST", "end time is required", http.StatusBadRequest)
		return
	}

	schedule, err := h.scheduleStorage.CreateSchedule(roomID, req.DaysOfWeek, req.StartTime, req.EndTime)
	if err != nil {
		msg := err.Error()
		switch {
		case msg == "room not found":
			domain.SendError(w, "ROOM_NOT_FOUND", "Room not found", http.StatusNotFound)
		case msg == "schedule already exists":
			domain.SendError(w, "SCHEDULE_EXISTS", "schedule for this room already exists and cannot be changed", http.StatusConflict)
		case strings.Contains(msg, "day of week must be between"):
			domain.SendError(w, "INVALID_REQUEST", msg, http.StatusBadRequest)
		case strings.Contains(msg, "duplicate day of week"):
			domain.SendError(w, "INVALID_REQUEST", msg, http.StatusBadRequest)
		case strings.Contains(msg, "invalid time format"):
			domain.SendError(w, "INVALID_REQUEST", msg, http.StatusBadRequest)
		case msg == "startTime must be before endTime":
			domain.SendError(w, "INVALID_REQUEST", "start time must be before endTime", http.StatusBadRequest)
		default:
			domain.SendError(w, "INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	domain.SendSuccess(w, map[string]interface{}{
		"schedule": schedule,
	}, http.StatusCreated)
}
