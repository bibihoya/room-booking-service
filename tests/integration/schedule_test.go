package integration

import (
	"testing"

	"github.com/internships-backend/test-backend-bibihoya-1/tests"
)

func TestCreateSchedule(t *testing.T) {
	tests.CleanDB()

	room, err := tests.RoomStorage.CreateRoom("Test Room", nil, nil)
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	schedule, err := tests.ScheduleStorage.CreateSchedule(room.ID, []int{1, 2, 3, 4, 5}, "09:00", "18:00")
	if err != nil {
		t.Fatalf("Failed to create schedule: %v", err)
	}

	if schedule.ID == "" {
		t.Error("Schedule ID should not be empty")
	}
	if schedule.RoomId != room.ID {
		t.Errorf("Expected roomId %s, got %s", room.ID, schedule.RoomId)
	}
	if schedule.StartTime != "09:00" {
		t.Errorf("Expected startTime 09:00, got %s", schedule.StartTime)
	}
	if schedule.EndTime != "18:00" {
		t.Errorf("Expected endTime 18:00, got %s", schedule.EndTime)
	}
}

func TestCreateScheduleDuplicate(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)

	_, err := tests.ScheduleStorage.CreateSchedule(room.ID, []int{1, 2, 3}, "09:00", "18:00")
	if err != nil {
		t.Fatalf("First schedule creation failed: %v", err)
	}

	_, err = tests.ScheduleStorage.CreateSchedule(room.ID, []int{1, 2, 3}, "09:00", "18:00")
	if err == nil {
		t.Error("Expected error for duplicate schedule, got nil")
	}
	if err.Error() != "schedule already exists" {
		t.Errorf("Expected 'schedule already exists', got %v", err)
	}
}

func TestCreateScheduleInvalidDays(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)

	_, err := tests.ScheduleStorage.CreateSchedule(room.ID, []int{8}, "09:00", "18:00")
	if err == nil {
		t.Error("Expected error for invalid day of week (8)")
	}

	_, err = tests.ScheduleStorage.CreateSchedule(room.ID, []int{1, 1, 2}, "09:00", "18:00")
	if err == nil {
		t.Error("Expected error for duplicate days")
	}
}

func TestCreateScheduleInvalidTime(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)

	_, err := tests.ScheduleStorage.CreateSchedule(room.ID, []int{1}, "29:00", "30:00")
	if err == nil {
		t.Error("Expected error for invalid time format")
	}

	_, err = tests.ScheduleStorage.CreateSchedule(room.ID, []int{1}, "18:00", "09:00")
	if err == nil {
		t.Error("Expected error for startTime after endTime")
	}
}

func TestCreateScheduleForNonExistentRoom(t *testing.T) {
	tests.CleanDB()

	_, err := tests.ScheduleStorage.CreateSchedule("00000000-0000-0000-0000-000000000000", []int{1}, "09:00", "18:00")
	if err == nil {
		t.Error("Expected error for non-existent room")
	}
}
