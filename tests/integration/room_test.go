package integration

import (
	"testing"

	"github.com/internships-backend/test-backend-bibihoya-1/tests"
)

func TestCreateRoom(t *testing.T) {
	tests.CleanDB()

	name := "Test Room"
	desc := "Test Description"
	cap := 10

	room, err := tests.RoomStorage.CreateRoom(name, &desc, &cap)
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	if room.ID == "" {
		t.Error("Room ID should not be empty")
	}
	if room.Name != name {
		t.Errorf("Expected name %s, got %s", name, room.Name)
	}
	if room.Description == nil || *room.Description != desc {
		t.Errorf("Expected description %s, got %v", desc, room.Description)
	}
	if room.Capacity == nil || *room.Capacity != cap {
		t.Errorf("Expected capacity %d, got %v", cap, room.Capacity)
	}
}

func TestListRooms(t *testing.T) {
	tests.CleanDB()

	room1, _ := tests.RoomStorage.CreateRoom("Room 1", nil, nil)
	room2, _ := tests.RoomStorage.CreateRoom("Room 2", nil, nil)

	rooms, err := tests.RoomStorage.RoomList()
	if err != nil {
		t.Fatalf("Failed to list rooms: %v", err)
	}

	if len(rooms) != 2 {
		t.Errorf("Expected 2 rooms, got %d", len(rooms))
	}

	found1, found2 := false, false
	for _, r := range rooms {
		if r.ID == room1.ID {
			found1 = true
		}
		if r.ID == room2.ID {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Error("Not all created rooms found in list")
	}
}
