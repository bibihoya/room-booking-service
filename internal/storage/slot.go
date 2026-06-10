package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/domain"
	"github.com/lib/pq"
)

type SlotStorage struct {
	db *sql.DB
}

func NewSlotStorage(db *sql.DB) *SlotStorage {
	return &SlotStorage{
		db: db,
	}
}

func (s *SlotStorage) GenerateSlotsForDate(roomID string, date time.Time) ([]domain.Slot, error) {
	weekday := int(date.Weekday())
	var pgDay int
	if weekday == 0 {
		pgDay = 7
	} else {
		pgDay = weekday
	}

	var daysOfWeek pq.Int64Array
	var startTime, endTime string

	query := `
			SELECT days_of_week, start_time, end_time
			FROM schedules
			WHERE room_id = $1 AND $2 = ANY(days_of_week)
	`

	err := s.db.QueryRow(query, roomID, pgDay).Scan(&daysOfWeek, &startTime, &endTime)
	if err == sql.ErrNoRows {
		return []domain.Slot{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get schedules: %w", err)
	}

	stHour, stMin := parseTime(startTime)
	endHour, endMin := parseTime(endTime)

	stDate := time.Date(date.Year(), date.Month(), date.Day(), stHour, stMin, 0, 0, time.UTC)
	endDate := time.Date(date.Year(), date.Month(), date.Day(), endHour, endMin, 0, 0, time.UTC)

	var slots []domain.Slot
	current := stDate

	for current.Add(30*time.Minute).Before(endDate) || current.Add(30*time.Minute).Equal(endDate) {
		slotEnd := current.Add(30 * time.Minute)

		var slotID string
		queryCheck := `
			SELECT id FROM slots
			WHERE room_id = $1 AND start_time = $2
	    `
		err := s.db.QueryRow(queryCheck, roomID, current).Scan(&slotID)
		if err == sql.ErrNoRows {
			slotID = uuid.New().String()
			queryIns := `
				INSERT INTO slots (id, room_id, start_time, end_time)
				VALUES ($1, $2, $3, $4)
			`
			_, err = s.db.Exec(queryIns, slotID, roomID, current, slotEnd)
			if err != nil {
				return nil, fmt.Errorf("failed to create slot: %w", err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("failed to check slot existence: %w", err)
		}

		slots = append(slots, domain.Slot{
			ID:     slotID,
			RoomId: roomID,
			Start:  current,
			End:    slotEnd,
		})

		current = slotEnd
	}

	return slots, nil
}

func (s *SlotStorage) GetAvailableSlots(roomID string, date time.Time) ([]domain.Slot, error) {
	slots, err := s.GenerateSlotsForDate(roomID, date)
	if err != nil {
		return nil, err
	}

	if len(slots) == 0 {
		return []domain.Slot{}, nil
	}

	slotIDs := make([]string, len(slots))
	for i, slot := range slots {
		slotIDs[i] = slot.ID
	}

	query := `
			SELECT DISTINCT slot_id
			FROM bookings
			WHERE slot_id = ANY($1) AND status = 'active'
	`

	rows, err := s.db.Query(query, pq.Array(slotIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to get booked slots: %w", err)
	}
	defer rows.Close()

	bookedMap := make(map[string]bool)
	for rows.Next() {
		var slotID string
		if err := rows.Scan(&slotID); err != nil {
			return nil, err
		}
		bookedMap[slotID] = true
	}

	var availableSlots []domain.Slot
	for _, slot := range slots {
		if !bookedMap[slot.ID] {
			availableSlots = append(availableSlots, slot)
		}
	}

	return availableSlots, nil
}

func parseTime(timeStr string) (int, int) {
	var hour, min int
	fmt.Sscanf(timeStr, "%d:%d", &hour, &min)
	return hour, min
}
