package integration

import (
	"testing"
	"time"

	"github.com/internships-backend/test-backend-bibihoya-1/tests"
)

func TestCreateBooking(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)
	_, err := tests.ScheduleStorage.CreateSchedule(room.ID, []int{5}, "09:00", "12:00")
	if err != nil {
		t.Fatalf("Failed to create schedule: %v", err)
	}

	date := time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC)
	slots, _ := tests.SlotStorage.GetAvailableSlots(room.ID, date)

	if len(slots) == 0 {
		t.Fatal("No slots available")
	}

	userID := "22222222-2222-2222-2222-222222222222"

	booking, err := tests.BookingStorage.CreateBooking(userID, slots[0].ID, false)
	if err != nil {
		t.Fatalf("Failed to create booking: %v", err)
	}

	if booking.ID == "" {
		t.Error("Booking ID should not be empty")
	}
	if booking.Status != "active" {
		t.Errorf("Expected status active, got %s", booking.Status)
	}
	if booking.UserId != userID {
		t.Errorf("Expected userID %s, got %s", userID, booking.UserId)
	}
}

func TestCreateBookingDuplicate(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)
	_, _ = tests.ScheduleStorage.CreateSchedule(room.ID, []int{5}, "09:00", "12:00")

	date := time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC)
	slots, _ := tests.SlotStorage.GetAvailableSlots(room.ID, date)

	userID := "22222222-2222-2222-2222-222222222222"

	_, err := tests.BookingStorage.CreateBooking(userID, slots[0].ID, false)
	if err != nil {
		t.Fatalf("First booking failed: %v", err)
	}

	_, err = tests.BookingStorage.CreateBooking(userID, slots[0].ID, false)
	if err == nil {
		t.Error("Expected error for duplicate booking")
	}
	if err.Error() != "slot already booked" {
		t.Errorf("Expected 'slot already booked', got %v", err)
	}
}

func TestCreateBookingPastSlot(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)
	_, _ = tests.ScheduleStorage.CreateSchedule(room.ID, []int{3}, "09:00", "12:00")

	date := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	slots, _ := tests.SlotStorage.GenerateSlotsForDate(room.ID, date)

	if len(slots) == 0 {
		t.Fatal("No slots generated")
	}

	userID := "22222222-2222-2222-2222-222222222222"

	_, err := tests.BookingStorage.CreateBooking(userID, slots[0].ID, false)
	if err == nil {
		t.Error("Expected error for past slot booking")
	}
	if err.Error() != "cannot book slot in the past" {
		t.Errorf("Expected 'cannot book slot in the past', got %v", err)
	}
}

func TestCancelBooking(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)
	_, _ = tests.ScheduleStorage.CreateSchedule(room.ID, []int{5}, "09:00", "12:00")

	date := time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC)
	slots, _ := tests.SlotStorage.GetAvailableSlots(room.ID, date)

	userID := "22222222-2222-2222-2222-222222222222"

	booking, _ := tests.BookingStorage.CreateBooking(userID, slots[0].ID, false)

	cancelled, err := tests.BookingStorage.CancelBooking(booking.ID, userID)
	if err != nil {
		t.Fatalf("Failed to cancel booking: %v", err)
	}

	if cancelled.Status != "cancelled" {
		t.Errorf("Expected status cancelled, got %s", cancelled.Status)
	}
}

func TestCancelBookingIdempotent(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)
	_, _ = tests.ScheduleStorage.CreateSchedule(room.ID, []int{5}, "09:00", "12:00")

	date := time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC)
	slots, _ := tests.SlotStorage.GetAvailableSlots(room.ID, date)

	userID := "22222222-2222-2222-2222-222222222222"

	booking, _ := tests.BookingStorage.CreateBooking(userID, slots[0].ID, false)

	_, err := tests.BookingStorage.CancelBooking(booking.ID, userID)
	if err != nil {
		t.Fatalf("First cancel failed: %v", err)
	}

	_, err = tests.BookingStorage.CancelBooking(booking.ID, userID)
	if err != nil {
		t.Fatalf("Second cancel (idempotent) failed: %v", err)
	}
}

func TestCancelBookingWrongUser(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)
	_, _ = tests.ScheduleStorage.CreateSchedule(room.ID, []int{5}, "09:00", "12:00")

	date := time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC)
	slots, _ := tests.SlotStorage.GetAvailableSlots(room.ID, date)

	userID := "22222222-2222-2222-2222-222222222222"
	wrongUserID := "11111111-1111-1111-1111-111111111111"

	booking, _ := tests.BookingStorage.CreateBooking(userID, slots[0].ID, false)

	_, err := tests.BookingStorage.CancelBooking(booking.ID, wrongUserID)
	if err == nil {
		t.Error("Expected error when cancelling another user's booking")
	}
	if err.Error() != "cannot cancel another user's booking" {
		t.Errorf("Expected 'cannot cancel another user's booking', got %v", err)
	}
}

func TestGetUserBookings(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)
	_, _ = tests.ScheduleStorage.CreateSchedule(room.ID, []int{5}, "09:00", "12:00")

	date := time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC)
	slots, _ := tests.SlotStorage.GetAvailableSlots(room.ID, date)

	userID := "22222222-2222-2222-2222-222222222222"

	tests.BookingStorage.CreateBooking(userID, slots[0].ID, false)
	tests.BookingStorage.CreateBooking(userID, slots[1].ID, false)

	bookings, err := tests.BookingStorage.GetUserBookings(userID)
	if err != nil {
		t.Fatalf("Failed to get user bookings: %v", err)
	}

	if len(bookings) != 2 {
		t.Errorf("Expected 2 bookings, got %d", len(bookings))
	}
}

func TestGetAllBookingsPaginated(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)
	_, _ = tests.ScheduleStorage.CreateSchedule(room.ID, []int{5}, "09:00", "12:00")

	date := time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC)
	slots, _ := tests.SlotStorage.GetAvailableSlots(room.ID, date)

	userID := "22222222-2222-2222-2222-222222222222"

	for i := 0; i < 5 && i < len(slots); i++ {
		tests.BookingStorage.CreateBooking(userID, slots[i].ID, false)
	}

	bookings, total, err := tests.BookingStorage.GetAllBookings(1, 3)
	if err != nil {
		t.Fatalf("Failed to get all bookings: %v", err)
	}

	if len(bookings) != 3 {
		t.Errorf("Expected 3 bookings on page 1, got %d", len(bookings))
	}
	if total != 5 {
		t.Errorf("Expected total 5 bookings, got %d", total)
	}

	bookings2, _, _ := tests.BookingStorage.GetAllBookings(2, 3)
	if len(bookings2) != 2 {
		t.Errorf("Expected 2 bookings on page 2, got %d", len(bookings2))
	}
}
