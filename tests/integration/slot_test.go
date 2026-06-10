package integration

import (
	"testing"
	"time"

	"github.com/internships-backend/test-backend-bibihoya-1/tests"
)

func TestGenerateSlots(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)
	_, err := tests.ScheduleStorage.CreateSchedule(room.ID, []int{5}, "09:00", "18:00")
	if err != nil {
		t.Fatalf("Failed to create schedule: %v", err)
	}

	date := time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC)

	slots, err := tests.SlotStorage.GenerateSlotsForDate(room.ID, date)
	if err != nil {
		t.Fatalf("Failed to generate slots: %v", err)
	}

	expectedCount := 18
	if len(slots) != expectedCount {
		t.Errorf("Expected %d slots, got %d", expectedCount, len(slots))
	}

	if len(slots) > 0 {
		expectedStart := time.Date(2026, 4, 3, 9, 0, 0, 0, time.UTC)
		if !slots[0].Start.Equal(expectedStart) {
			t.Errorf("Expected start %v, got %v", expectedStart, slots[0].Start)
		}
	}
}

func TestGetAvailableSlots(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)
	_, err := tests.ScheduleStorage.CreateSchedule(room.ID, []int{5}, "09:00", "12:00")
	if err != nil {
		t.Fatalf("Failed to create schedule: %v", err)
	}

	date := time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC)

	slots, err := tests.SlotStorage.GetAvailableSlots(room.ID, date)
	if err != nil {
		t.Fatalf("Failed to get available slots: %v", err)
	}

	expectedCount := 6
	if len(slots) != expectedCount {
		t.Errorf("Expected %d available slots, got %d", expectedCount, len(slots))
	}
}

func TestGetAvailableSlotsNoSchedule(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)

	date := time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC)

	slots, err := tests.SlotStorage.GetAvailableSlots(room.ID, date)
	if err != nil {
		t.Fatalf("Failed to get available slots: %v", err)
	}

	if len(slots) != 0 {
		t.Errorf("Expected 0 slots for room without schedule, got %d", len(slots))
	}
}

func TestGetAvailableSlotsWrongDay(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)
	_, err := tests.ScheduleStorage.CreateSchedule(room.ID, []int{5}, "09:00", "18:00")
	if err != nil {
		t.Fatalf("Failed to create schedule: %v", err)
	}

	date := time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC)

	slots, err := tests.SlotStorage.GetAvailableSlots(room.ID, date)
	if err != nil {
		t.Fatalf("Failed to get available slots: %v", err)
	}

	if len(slots) != 0 {
		t.Errorf("Expected 0 slots for wrong day, got %d", len(slots))
	}
}
