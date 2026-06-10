package storage

import (
	"database/sql"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/domain"
	"github.com/lib/pq"
)

type ScheduleStorage struct {
	db *sql.DB
}

func NewScheduleStorage(db *sql.DB) *ScheduleStorage {
	return &ScheduleStorage{db: db}
}

func validateTime(timeStr string) error {
	matched, err := regexp.MatchString(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`, timeStr)
	if err != nil {
		return fmt.Errorf("failed to validate time: %w", err)
	}
	if !matched {
		return fmt.Errorf("invalid time format, expected HH:MM, got %s", timeStr)
	}
	return nil
}

func compareTimes(start, end string) (bool, error) {
	startParts := strings.Split(start, ":")
	endParts := strings.Split(end, ":")

	stHour, _ := strconv.Atoi(startParts[0])
	stMin, _ := strconv.Atoi(startParts[1])
	endHour, _ := strconv.Atoi(endParts[0])
	endMin, _ := strconv.Atoi(endParts[1])

	stTotal := stHour*60 + stMin
	endTotal := endHour*60 + endMin

	return stTotal < endTotal, nil
}

func (s *ScheduleStorage) CreateSchedule(roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
	var exists bool
	queryEx := `
			SELECT EXISTS (SELECT 1 FROM rooms WHERE id = $1)
	`
	err := s.db.QueryRow(queryEx, roomID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failedto check room existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("room does not exist")
	}

	var schedExists bool
	queryScEx := `
			SELECT EXISTS (SELECT 1 FROM schedules WHERE room_id = $1)
	`
	err = s.db.QueryRow(queryScEx, roomID).Scan(&schedExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check schedule existence: %w", err)
	}
	if schedExists {
		return nil, fmt.Errorf("schedule already exists")
	}

	if err := validateTime(startTime); err != nil {
		return nil, err
	}
	if err := validateTime(endTime); err != nil {
		return nil, err
	}

	valid, err := compareTimes(startTime, endTime)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("startTime must be before endTime")
	}

	uniqueDays := make(map[int]bool)

	for _, day := range daysOfWeek {
		if day < 1 || day > 7 {
			return nil, fmt.Errorf("day of week must be between 1 and 7, got %d", day)
		}
		if uniqueDays[day] {
			return nil, fmt.Errorf("duplicate day of week: %d", day)
		}
		uniqueDays[day] = true
	}
	uniqueDaysSlice := make([]int, 0, len(uniqueDays))
	for day := range uniqueDays {
		uniqueDaysSlice = append(uniqueDaysSlice, day)
	}

	sort.Ints(uniqueDaysSlice)

	id := uuid.New().String()
	queryIns := `
			INSERT INTO schedules (id, room_id, days_of_week, start_time, end_time)
			VALUES ($1, $2, $3, $4, $5)
	`
	_, err = s.db.Exec(queryIns, id, roomID, pq.Array(uniqueDaysSlice), startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to insert schedule: %w", err)
	}

	return &domain.Schedule{
		ID:         id,
		RoomId:     roomID,
		DaysOfWeek: uniqueDaysSlice,
		StartTime:  startTime,
		EndTime:    endTime,
	}, nil
}
