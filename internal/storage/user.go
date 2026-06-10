package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) CreateUser(email, password, role string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	id := uuid.New().String()

	query := `
			INSERT INTO users (id, email, password_hash, role)
        	VALUES ($1, $2, $3, $4)
	`

	_, err = s.db.Exec(query, id, email, string(hashedPassword), role)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &domain.User{
		ID:        id,
		Email:     email,
		Role:      role,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (s *UserStorage) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	var passwordHash string

	query := `
			SELECT id, email, password_hash, role, created_at
        	FROM users WHERE email = $1
	`

	err := s.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &passwordHash, &user.Role, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.PasswordHash = passwordHash
	return &user, nil
}

func (s *UserStorage) CheckPassword(user *domain.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}
