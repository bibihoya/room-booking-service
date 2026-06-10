package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}
type Room struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Capacity    *int      `json:"capacity,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Schedule struct {
	ID         string `json:"id"`
	RoomId     string `json:"roomId"`
	DaysOfWeek []int  `json:"daysOfWeek"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

type Slot struct {
	ID     string    `json:"id"`
	RoomId string    `json:"roomId"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}

type Booking struct {
	ID             string    `json:"id"`
	SlotId         string    `json:"slotId"`
	UserId         string    `json:"userId"`
	Status         string    `json:"status"`
	ConferenceLink *string   `json:"conferenceLink,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}
