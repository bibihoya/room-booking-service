package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/domain"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/middleware"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/storage"
)

type BookingHandler struct {
	bookingStorage *storage.BookingStorage
}

func NewBookingHandler(bookingStorage *storage.BookingStorage) *BookingHandler {
	return &BookingHandler{bookingStorage: bookingStorage}
}

// CreateBooking godoc
// @Summary      Создать бронь
// @Description  Создаёт бронь на слот. Только user. userId берётся из JWT.
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object true "Данные брони" example({"slotId": "uuid", "createConferenceLink": false})
// @Success      201  {object}  map[string]interface{}  "booking"
// @Failure      400  {object}  map[string]interface{}  "error"
// @Failure      401  {object}  map[string]interface{}  "error"
// @Failure      403  {object}  map[string]interface{}  "error"
// @Failure      404  {object}  map[string]interface{}  "error"
// @Failure      409  {object}  map[string]interface{}  "error"
// @Failure      500  {object}  map[string]interface{}  "error"
// @Router       /bookings/create [post]
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		domain.SendError(w, "UNAUTHORIZED", "User not found in context", http.StatusUnauthorized)
		return
	}

	var req struct {
		SlotID               string `json:"slotId"`
		CreateConferenceLink bool   `json:"createConferenceLink"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.SendError(w, "INVALID_REQUEST", "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SlotID == "" {
		domain.SendError(w, "INVALID_REQUEST", "slotId is required", http.StatusBadRequest)
		return
	}

	booking, err := h.bookingStorage.CreateBooking(userID, req.SlotID, req.CreateConferenceLink)
	if err != nil {
		switch err.Error() {
		case "slot not found":
			domain.SendError(w, "SLOT_NOT_FOUND", "Slot not found", http.StatusNotFound)
		case "cannot book slot in the past":
			domain.SendError(w, "INVALID_REQUEST", "Cannot book slot in the past", http.StatusBadRequest)
		case "slot already booked":
			domain.SendError(w, "SLOT_ALREADY_BOOKED", "Slot already booked", http.StatusConflict)
		default:
			domain.SendError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
		}
		return
	}
	domain.SendSuccess(w, map[string]interface{}{
		"booking": booking,
	}, http.StatusCreated)
}

// CancelBooking godoc
// @Summary      Отменить бронь
// @Description  Отменяет бронь. Идемпотентно. Только владелец брони.
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        bookingId path string true "ID брони"
// @Success      200  {object}  map[string]interface{}  "booking"
// @Failure      401  {object}  map[string]interface{}  "error"
// @Failure      403  {object}  map[string]interface{}  "error"
// @Failure      404  {object}  map[string]interface{}  "error"
// @Failure      500  {object}  map[string]interface{}  "error"
// @Router       /bookings/{bookingId}/cancel [post]
func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		domain.SendError(w, "UNAUTHORIZED", "User not found in context", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	bookingID := vars["bookingId"]

	booking, err := h.bookingStorage.CancelBooking(bookingID, userID)
	if err != nil {
		switch err.Error() {
		case "booking not found":
			domain.SendError(w, "BOOKING_NOT_FOUND", "Booking not found", http.StatusNotFound)
		case "cannot cancel another user's booking":
			domain.SendError(w, "FORBIDDEN", "Cannot cancel another user's booking", http.StatusForbidden)
		default:
			domain.SendError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
		}
		return
	}

	domain.SendSuccess(w, map[string]interface{}{
		"booking": booking,
	}, http.StatusOK)
}

// GetMyBooking godoc
// @Summary      Мои брони
// @Description  Возвращает список активных броней текущего пользователя на будущие слоты.
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "bookings"
// @Failure      401  {object}  map[string]interface{}  "error"
// @Failure      403  {object}  map[string]interface{}  "error"
// @Failure      500  {object}  map[string]interface{}  "error"
// @Router       /bookings/my [get]
func (h *BookingHandler) GetMyBooking(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		domain.SendError(w, "UNAUTHORIZED", "User not found in context", http.StatusUnauthorized)
		return
	}

	bookings, err := h.bookingStorage.GetUserBookings(userID)
	if err != nil {
		domain.SendError(w, "INTERNAL_ERROR", "Failed to get bookings", http.StatusInternalServerError)
		return
	}

	domain.SendSuccess(w, map[string]interface{}{
		"bookings": bookings,
	}, http.StatusOK)
}

// GetAllBookings godoc
// @Summary      Список всех броней (admin)
// @Description  Возвращает список всех броней с пагинацией. Только admin.
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Номер страницы" default(1) minimum(1)
// @Param        pageSize query int false "Размер страницы" default(20) minimum(1) maximum(100)
// @Success      200  {object}  map[string]interface{}  "bookings, pagination"
// @Failure      400  {object}  map[string]interface{}  "error"
// @Failure      401  {object}  map[string]interface{}  "error"
// @Failure      403  {object}  map[string]interface{}  "error"
// @Failure      500  {object}  map[string]interface{}  "error"
// @Router       /bookings/list [get]
func (h *BookingHandler) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := 1
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			domain.SendError(w, "INVALID_REQUEST", "page must be a positive integer", http.StatusBadRequest)
			return
		}
		page = p
	}

	pageSize := 20
	if pageSizeStr != "" {
		ps, err := strconv.Atoi(pageSizeStr)
		if err != nil || ps < 1 {
			domain.SendError(w, "INVALID_REQUEST", "pageSize must be a positive integer", http.StatusBadRequest)
			return
		}
		if ps > 100 {
			domain.SendError(w, "INVALID_REQUEST", "pageSize must be <= 100", http.StatusBadRequest)
			return
		}
		pageSize = ps
	}

	bookings, total, err := h.bookingStorage.GetAllBookings(page, pageSize)
	if err != nil {
		domain.SendError(w, "INTERNAL_ERROR", "Failed to get bookings", http.StatusInternalServerError)
		return
	}

	domain.SendSuccess(w, map[string]interface{}{
		"bookings": bookings,
		"pagination": map[string]interface{}{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
		},
	}, http.StatusOK)
}
