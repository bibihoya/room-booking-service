package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/domain"
)

type RoomStorage struct {
	db *sql.DB
}

func NewRoomStorage(db *sql.DB) *RoomStorage {
	return &RoomStorage{db: db}
}

func (s *RoomStorage) RoomList() ([]domain.Room, error) {
	query := `
			SELECT id, name, description, capacity, created_at
			FROM rooms
			ORDER BY created_at DESC
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []domain.Room
	for rows.Next() {
		var room domain.Room
		var description sql.NullString
		var capacity sql.NullInt32

		if err := rows.Scan(&room.ID, &room.Name, &description, &capacity, &room.CreatedAt); err != nil {
			return nil, err
		}

		if description.Valid {
			room.Description = &description.String
		}
		if capacity.Valid {
			cap := int(capacity.Int32)
			room.Capacity = &cap
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (s *RoomStorage) CreateRoom(name string, description *string, capacity *int) (*domain.Room, error) {
	id := uuid.New().String()
	now := time.Now().UTC()

	query := `
			INSERT INTO rooms (id, name, description, capacity, created_at)
			VALUES ($1, $2, $3, $4, $5)
	`
	_, err := s.db.Exec(query, id, name, description, capacity, now)
	if err != nil {
		return nil, err
	}

	return &domain.Room{
		ID:          id,
		Name:        name,
		Description: description,
		Capacity:    capacity,
		CreatedAt:   now,
	}, nil
}
