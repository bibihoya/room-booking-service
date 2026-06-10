package storage

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/domain"
)

type BookingStorage struct {
	db *sql.DB
}

func NewBookingStorage(db *sql.DB) *BookingStorage {
	return &BookingStorage{db: db}
}

// мок запроса во внешний сервис
func CreateMeeting(bookingID string) (string, error) {
	time.Sleep(100 * time.Millisecond)

	// Допустим в 5 процентах случаев ошибка сервиса
	if rand.Intn(100) < 5 {
		return "", fmt.Errorf("conference service unavailable: 500 Internal Server Error")
	}

	// Здесь должен быть HTTP запрос:
	// resp, err := http.Post("https://api.conference.com/meetings", ...)
	conferenceURL := fmt.Sprintf("https://meet.conference-service.com/room/%s", bookingID)

	return conferenceURL, nil
}

func (s *BookingStorage) CreateBooking(userID, slotID string, createConferenceLink bool) (*domain.Booking, error) {
	var slotStart time.Time
	var slotRoomID string
	query := `
			SELECT start_time, room_id FROM slots
			WHERE id = $1
	`
	err := s.db.QueryRow(query, slotID).Scan(&slotStart, &slotRoomID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("slot not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to check slot: %w", err)
	}

	if slotStart.Before(time.Now().UTC()) {
		return nil, fmt.Errorf("cannot book slot in the past")
	}

	var existBookID string
	query = `
			SELECT id FROM bookings
			WHERE slot_id = $1 AND status = 'active'
	`
	err = s.db.QueryRow(query, slotID).Scan(&existBookID)
	if err == nil {
		return nil, fmt.Errorf("slot already booked")
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check existing booking: %w", err)
	}

	bookingID := uuid.New().String()
	now := time.Now().UTC()

	var conferenceLink *string
	if createConferenceLink {
		link, err := CreateMeeting(bookingID)
		if err != nil {
			log.Printf("failed to generate link: %v", err)
		} else {
			conferenceLink = &link
		}
	}

	queryIns := `
			INSERT INTO bookings (id, slot_id, user_id, status, conference_link, created_at, updated_at)
			VALUES ($1, $2, $3, 'active', $4, $5, $5)
	`
	_, err = s.db.Exec(queryIns, bookingID, slotID, userID, conferenceLink, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	return &domain.Booking{
		ID:             bookingID,
		SlotId:         slotID,
		UserId:         userID,
		Status:         "active",
		ConferenceLink: conferenceLink,
		CreatedAt:      now,
	}, nil
}

func (s *BookingStorage) CancelBooking(bookingID, userID string) (*domain.Booking, error) {
	var exStatus, exUserID, slotID string
	var createdAt time.Time
	var conferenceLink sql.NullString
	query := `
			SELECT status, user_id, slot_id, created_at, conference_link
			FROM bookings
			WHERE id = $1
	`
	err := s.db.QueryRow(query, bookingID).Scan(&exStatus, &exUserID, &slotID, &createdAt, &conferenceLink)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("booking not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to check booking: %w", err)
	}

	if exUserID != userID {
		return nil, fmt.Errorf("cannot cancel another user's booking")
	}

	if exStatus == "cancelled" {
		var link *string
		if conferenceLink.Valid {
			link = &conferenceLink.String
		}

		return &domain.Booking{
			ID:             bookingID,
			SlotId:         slotID,
			UserId:         userID,
			Status:         "cancelled",
			ConferenceLink: link,
			CreatedAt:      createdAt,
		}, nil
	}

	queryUpd := `
			UPDATE bookings
			SET status = 'cancelled', updated_at = NOW()
			WHERE id = $1
	`
	_, err = s.db.Exec(queryUpd, bookingID)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel booking: %w", err)
	}

	var link *string
	if conferenceLink.Valid {
		link = &conferenceLink.String
	}
	return &domain.Booking{
		ID:             bookingID,
		SlotId:         slotID,
		UserId:         userID,
		Status:         "cancelled",
		ConferenceLink: link,
		CreatedAt:      createdAt,
	}, nil
}

func (s *BookingStorage) GetUserBookings(userID string) ([]domain.Booking, error) {
	query := `
			SELECT b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at, s.start_time
			FROM bookings b
			JOIN slots s ON b.slot_id = s.id
			WHERE b.user_id = $1 AND b.status = 'active' AND s.start_time >= NOW()
			ORDER BY s.start_time ASC
	`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user bookings: %w", err)
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var booking domain.Booking
		var conferenceLink sql.NullString
		var slotStart time.Time

		if err := rows.Scan(&booking.ID, &booking.SlotId, &booking.UserId, &booking.Status, &conferenceLink, &booking.CreatedAt, &slotStart); err != nil {
			return nil, err
		}

		if conferenceLink.Valid {
			booking.ConferenceLink = &conferenceLink.String
		}

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (s *BookingStorage) GetAllBookings(page, pageSize int) ([]domain.Booking, int, error) {
	var total int
	query := `
			SELECT COUNT(*) FROM bookings
	`
	err := s.db.QueryRow(query).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	offset := (page - 1) * pageSize

	query = `
			SELECT id, slot_id, user_id, status, conference_link, created_at
			FROM bookings
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
	`
	rows, err := s.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get bookings: %w", err)
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var booking domain.Booking
		var conferenceLink sql.NullString

		if err := rows.Scan(&booking.ID, &booking.SlotId, &booking.UserId, &booking.Status, &conferenceLink, &booking.CreatedAt); err != nil {
			return nil, 0, err
		}

		if conferenceLink.Valid {
			booking.ConferenceLink = &conferenceLink.String
		}

		bookings = append(bookings, booking)
	}

	return bookings, total, nil
}
