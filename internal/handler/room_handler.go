package handler

import (
	"encoding/json"
	"net/http"

	"github.com/internships-backend/test-backend-bibihoya-1/internal/domain"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/storage"
)

type RoomHandler struct {
	roomStorage *storage.RoomStorage
}

func NewRoomHandler(store *storage.RoomStorage) *RoomHandler {
	return &RoomHandler{
		roomStorage: store,
	}
}

// RoomList godoc
// @Summary      Список переговорок
// @Description  Возвращает список всех переговорок. Доступно admin и user.
// @Tags         Rooms
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "rooms"
// @Failure      401  {object}  map[string]interface{}  "error"
// @Failure      500  {object}  map[string]interface{}  "error"
// @Router       /rooms/list [get]
func (h *RoomHandler) RoomList(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.roomStorage.RoomList()
	if err != nil {
		domain.SendError(w, "INTERNAL_ERROR", "Failed to list rooms", http.StatusInternalServerError)
		return
	}

	domain.SendSuccess(w, map[string]interface{}{
		"rooms": rooms,
	}, http.StatusOK)
}

// RoomCreate godoc
// @Summary      Создать переговорку
// @Description  Создаёт новую переговорку. Только для admin.
// @Tags         Rooms
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object true "Данные комнаты" example({"name": "Conference Room", "description": "Big room", "capacity": 10})
// @Success      201  {object}  map[string]interface{}  "room"
// @Failure      400  {object}  map[string]interface{}  "error"
// @Failure      401  {object}  map[string]interface{}  "error"
// @Failure      403  {object}  map[string]interface{}  "error"
// @Failure      500  {object}  map[string]interface{}  "error"
// @Router       /rooms/create [post]
func (h *RoomHandler) RoomCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
		Capacity    *int    `json:"capacity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.SendError(w, "INVALID_REQUEST", "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		domain.SendError(w, "INVALID_REQUEST", "Name is required", http.StatusBadRequest)
		return
	}

	room, err := h.roomStorage.CreateRoom(req.Name, req.Description, req.Capacity)
	if err != nil {
		domain.SendError(w, "INTERNAL_ERROR", "Failed to create room", http.StatusInternalServerError)
		return
	}

	domain.SendSuccess(w, map[string]interface{}{
		"room": room,
	}, http.StatusCreated)
}
