package e2e

import (
	"testing"
	"time"

	"github.com/internships-backend/test-backend-bibihoya-1/tests"
)

func TestFullBookingFlow(t *testing.T) {
	tests.CleanDB()

	room, err := tests.RoomStorage.CreateRoom("Conference Room", nil, nil)
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}
	t.Logf("Room created: %s", room.ID)

	_, err = tests.ScheduleStorage.CreateSchedule(room.ID, []int{1, 2, 3, 4, 5}, "09:00", "18:00")
	if err != nil {
		t.Fatalf("Failed to create schedule: %v", err)
	}
	t.Log("Schedule created")

	date := time.Now().AddDate(0, 0, 7)
	for date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		date = date.AddDate(0, 0, 1)
	}

	slots, err := tests.SlotStorage.GetAvailableSlots(room.ID, date)
	if err != nil {
		t.Fatalf("Failed to get slots: %v", err)
	}
	if len(slots) == 0 {
		t.Fatal("No slots available")
	}
	t.Logf("Got %d available slots", len(slots))

	userID := "22222222-2222-2222-2222-222222222222"
	booking, err := tests.BookingStorage.CreateBooking(userID, slots[0].ID, false)
	if err != nil {
		t.Fatalf("Failed to create booking: %v", err)
	}
	t.Logf("Booking created: %s", booking.ID)

	availableSlots, _ := tests.SlotStorage.GetAvailableSlots(room.ID, date)
	for _, slot := range availableSlots {
		if slot.ID == slots[0].ID {
			t.Error("Booked slot is still available")
		}
	}
	t.Log("Booked slot no longer available")

	cancelled, err := tests.BookingStorage.CancelBooking(booking.ID, userID)
	if err != nil {
		t.Fatalf("Failed to cancel booking: %v", err)
	}
	if cancelled.Status != "cancelled" {
		t.Errorf("Expected status cancelled, got %s", cancelled.Status)
	}
	t.Log("Booking cancelled")

	t.Log("Full booking flow completed successfully!")
}

func TestAdminGetAllBookings(t *testing.T) {
	tests.CleanDB()

	room, _ := tests.RoomStorage.CreateRoom("Test Room", nil, nil)
	_, _ = tests.ScheduleStorage.CreateSchedule(room.ID, []int{5}, "09:00", "12:00")

	date := time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC)
	slots, _ := tests.SlotStorage.GetAvailableSlots(room.ID, date)

	userID := "22222222-2222-2222-2222-222222222222"
	for i := 0; i < 3 && i < len(slots); i++ {
		tests.BookingStorage.CreateBooking(userID, slots[i].ID, false)
	}

	bookings, total, err := tests.BookingStorage.GetAllBookings(1, 10)
	if err != nil {
		t.Fatalf("Failed to get all bookings: %v", err)
	}

	if total != 3 {
		t.Errorf("Expected total 3 bookings, got %d", total)
	}
	if len(bookings) != 3 {
		t.Errorf("Expected 3 bookings, got %d", len(bookings))
	}
	t.Logf("Admin retrieved %d bookings", total)
}
