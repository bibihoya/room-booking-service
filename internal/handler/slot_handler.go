package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/domain"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/storage"
)

type SlotHandler struct {
	slotStorage *storage.SlotStorage
	roomStorage *storage.RoomStorage
}

func NewSlotHandler(slotStorage *storage.SlotStorage, roomStorage *storage.RoomStorage) *SlotHandler {
	return &SlotHandler{
		slotStorage: slotStorage,
		roomStorage: roomStorage,
	}
}

// ListAvailableSlots godoc
// @Summary      Список доступных слотов
// @Description  Возвращает слоты, не занятые активной бронью, для указанной переговорки на указанную дату.
// @Tags         Slots
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        roomId path string true "ID переговорки"
// @Param        date query string true "Дата в формате YYYY-MM-DD" example(2026-04-10)
// @Success      200  {object}  map[string]interface{}  "slots"
// @Failure      400  {object}  map[string]interface{}  "error"
// @Failure      401  {object}  map[string]interface{}  "error"
// @Failure      404  {object}  map[string]interface{}  "error"
// @Failure      500  {object}  map[string]interface{}  "error"
// @Router       /rooms/{roomId}/slots/list [get]
func (h *SlotHandler) ListAvailableSlots(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomId"]

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		domain.SendError(w, "INVALID_REQUEST", "date is required", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		domain.SendError(w, "INVALID_REQUEST", "invalid date format, expected UTC", http.StatusBadRequest)
		return
	}

	rooms, err := h.roomStorage.RoomList()
	if err != nil {
		domain.SendError(w, "INTERNAL_ERROR", "failed to check room", http.StatusInternalServerError)
		return
	}

	exists := false
	for _, room := range rooms {
		if room.ID == roomID {
			exists = true
			break
		}
	}

	if !exists {
		domain.SendError(w, "ROOM_NOT_FOUND", "room not found", http.StatusNotFound)
		return
	}

	slots, err := h.slotStorage.GetAvailableSlots(roomID, date)
	if err != nil {
		domain.SendError(w, "INTERNAL_ERROR", "failed to get available slots", http.StatusInternalServerError)
		return
	}

	domain.SendSuccess(w, map[string]interface{}{
		"slots": slots,
	}, http.StatusOK)
}
