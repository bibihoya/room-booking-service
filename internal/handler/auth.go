package handler

import (
	"encoding/json"
	"net/http"

	"github.com/internships-backend/test-backend-bibihoya-1/internal/domain"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/storage"
	"github.com/internships-backend/test-backend-bibihoya-1/pkg/jwt"
)

type AuthHandler struct {
	jwtSecret   string
	jwtTTL      int
	userStorage *storage.UserStorage
}

func NewAuthHandler(jwtSecret string, jwtTTL int, userStorage *storage.UserStorage) *AuthHandler {
	return &AuthHandler{
		jwtSecret:   jwtSecret,
		jwtTTL:      jwtTTL,
		userStorage: userStorage,
	}
}

// DummyLogin godoc
// @Summary      Получить тестовый JWT
// @Description  Выдаёт тестовый JWT для указанной роли (admin/user). Для каждой роли возвращается фиксированный UUID.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body object true "Роль пользователя" example({"role": "admin"})
// @Success      200  {object}  map[string]interface{}  "token"
// @Failure      400  {object}  map[string]interface{}  "error"
// @Failure      500  {object}  map[string]interface{}  "error"
// @Router       /dummyLogin [post]
func (h *AuthHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Role string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.SendError(w, "INVALID_REQUEST", "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Role != "admin" && req.Role != "user" {
		domain.SendError(w, "INVALID_REQUEST", "Role must be admin or user", http.StatusBadRequest)
		return
	}

	var userID string
	if req.Role == "admin" {
		userID = "11111111-1111-1111-1111-111111111111"
	} else {
		userID = "22222222-2222-2222-2222-222222222222"
	}

	token, err := jwt.GenerateToken(userID, req.Role, h.jwtSecret, h.jwtTTL)
	if err != nil {
		domain.SendError(w, "INTERNAL_ERROR", "Failed to generate token", http.StatusInternalServerError)
		return
	}

	domain.SendSuccess(w, map[string]interface{}{
		"token": token,
	}, http.StatusOK)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.SendError(w, "INVALID_REQUEST", "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		domain.SendError(w, "INVALID_REQUEST", "Email is required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		domain.SendError(w, "INVALID_REQUEST", "Password is required", http.StatusBadRequest)
		return
	}

	if req.Role != "admin" && req.Role != "user" {
		req.Role = "user"
	}

	user, err := h.userStorage.CreateUser(req.Email, req.Password, req.Role)
	if err != nil {
		domain.SendError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
		return
	}

	domain.SendSuccess(w, map[string]interface{}{
		"user": user,
	}, http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.SendError(w, "INVALID_REQUEST", "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		domain.SendError(w, "INVALID_REQUEST", "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.userStorage.GetUserByEmail(req.Email)
	if err != nil {
		domain.SendError(w, "UNAUTHORIZED", "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !h.userStorage.CheckPassword(user, req.Password) {
		domain.SendError(w, "UNAUTHORIZED", "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Role, h.jwtSecret, h.jwtTTL)
	if err != nil {
		domain.SendError(w, "INTERNAL_ERROR", "Failed to generate token", http.StatusInternalServerError)
		return
	}

	domain.SendSuccess(w, map[string]string{"token": token}, http.StatusOK)
}
