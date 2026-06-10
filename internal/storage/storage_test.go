package storage

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewRoomStorage(t *testing.T) {
	db := &sql.DB{}
	storage := NewRoomStorage(db)
	assert.NotNil(t, storage)
	assert.Equal(t, db, storage.db)
}

func TestRoomStorage_CreateRoom(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	storage := NewRoomStorage(db)
	name := "Test Room"
	desc := "Description"
	cap := 10

	mock.ExpectExec("INSERT INTO rooms").
		WithArgs(sqlmock.AnyArg(), name, desc, cap, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	room, err := storage.CreateRoom(name, &desc, &cap)
	assert.NoError(t, err)
	assert.NotEmpty(t, room.ID)
	assert.Equal(t, name, room.Name)
}

func TestRoomStorage_CreateRoomWithoutOptional(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	storage := NewRoomStorage(db)
	name := "Simple Room"

	mock.ExpectExec("INSERT INTO rooms").
		WithArgs(sqlmock.AnyArg(), name, nil, nil, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	room, err := storage.CreateRoom(name, nil, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, room.ID)
	assert.Equal(t, name, room.Name)
	assert.Nil(t, room.Description)
	assert.Nil(t, room.Capacity)
}

func TestRoomStorage_ListRooms(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	storage := NewRoomStorage(db)
	now := time.Now().UTC()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "capacity", "created_at"}).
		AddRow(uuid.New().String(), "Room 1", "Desc 1", 10, now).
		AddRow(uuid.New().String(), "Room 2", nil, nil, now)

	mock.ExpectQuery("SELECT id, name, description, capacity, created_at FROM rooms").
		WillReturnRows(rows)

	rooms, err := storage.RoomList()
	assert.NoError(t, err)
	assert.Len(t, rooms, 2)
}

func TestRoomStorage_ListRoomsEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	storage := NewRoomStorage(db)
	rows := sqlmock.NewRows([]string{"id", "name", "description", "capacity", "created_at"})

	mock.ExpectQuery("SELECT id, name, description, capacity, created_at FROM rooms").
		WillReturnRows(rows)

	rooms, err := storage.RoomList()
	assert.NoError(t, err)
	assert.Len(t, rooms, 0)
}

func TestNewScheduleStorage(t *testing.T) {
	db := &sql.DB{}
	storage := NewScheduleStorage(db)
	assert.NotNil(t, storage)
	assert.Equal(t, db, storage.db)
}

func TestValidateTime(t *testing.T) {
	tests := []struct {
		name    string
		timeStr string
		wantErr bool
	}{
		{"valid 09:00", "09:00", false},
		{"valid 9:00", "9:00", false},
		{"valid 18:30", "18:30", false},
		{"valid 00:00", "00:00", false},
		{"invalid 24:00", "24:00", true},
		{"invalid 09:60", "09:60", true},
		{"invalid format", "9-00", true},
		{"empty", "", true},
		{"too long", "09:000", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTime(tt.timeStr)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCompareTimes(t *testing.T) {
	tests := []struct {
		name  string
		start string
		end   string
		want  bool
	}{
		{"start before end", "09:00", "18:00", true},
		{"start equal end", "09:00", "09:00", false},
		{"start after end", "18:00", "09:00", false},
		{"different minutes", "09:30", "10:00", true},
		{"same hour diff minutes", "09:30", "09:45", true},
		{"late night", "23:00", "23:30", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := compareTimes(tt.start, tt.end)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestScheduleStorage_CreateSchedule(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	storage := NewScheduleStorage(db)
	roomID := "11111111-1111-1111-1111-111111111111"
	daysOfWeek := []int{1, 2, 3}
	startTime := "09:00"
	endTime := "18:00"

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(roomID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(roomID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec("INSERT INTO schedules").
		WithArgs(sqlmock.AnyArg(), roomID, sqlmock.AnyArg(), startTime, endTime).
		WillReturnResult(sqlmock.NewResult(1, 1))

	schedule, err := storage.CreateSchedule(roomID, daysOfWeek, startTime, endTime)
	assert.NoError(t, err)
	assert.NotEmpty(t, schedule.ID)
	assert.Equal(t, roomID, schedule.RoomId)
}

func TestNewSlotStorage(t *testing.T) {
	db := &sql.DB{}
	storage := NewSlotStorage(db)
	assert.NotNil(t, storage)
	assert.Equal(t, db, storage.db)
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		name     string
		timeStr  string
		wantHour int
		wantMin  int
	}{
		{"standard format", "09:30", 9, 30},
		{"without leading zero", "9:30", 9, 30},
		{"evening time", "18:00", 18, 0},
		{"with minutes", "14:45", 14, 45},
		{"midnight", "00:00", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hour, min := parseTime(tt.timeStr)
			assert.Equal(t, tt.wantHour, hour)
			assert.Equal(t, tt.wantMin, min)
		})
	}
}

func TestWeekdayConversion(t *testing.T) {
	tests := []struct {
		goWeekday  time.Weekday
		expectedPG int
	}{
		{time.Monday, 1},
		{time.Tuesday, 2},
		{time.Wednesday, 3},
		{time.Thursday, 4},
		{time.Friday, 5},
		{time.Saturday, 6},
		{time.Sunday, 7},
	}

	for _, tt := range tests {
		t.Run(tt.goWeekday.String(), func(t *testing.T) {
			weekday := int(tt.goWeekday)
			var pgDay int
			if weekday == 0 {
				pgDay = 7
			} else {
				pgDay = weekday
			}
			assert.Equal(t, tt.expectedPG, pgDay)
		})
	}
}

func TestNewBookingStorage(t *testing.T) {
	db := &sql.DB{}
	storage := NewBookingStorage(db)
	assert.NotNil(t, storage)
	assert.Equal(t, db, storage.db)
}

func TestValidationIntegration(t *testing.T) {
	validTimes := []string{"09:00", "18:30", "00:00", "23:59"}
	for _, start := range validTimes {
		for _, end := range validTimes {
			err := validateTime(start)
			if err != nil {
				continue
			}
			err = validateTime(end)
			if err != nil {
				continue
			}
			_, err = compareTimes(start, end)
			assert.NoError(t, err)
		}
	}
}

func TestRoomStorage_CreateRoom_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	storage := NewRoomStorage(db)

	mock.ExpectExec("INSERT INTO rooms").
		WillReturnError(sql.ErrConnDone)

	_, err = storage.CreateRoom("Test", nil, nil)
	assert.Error(t, err)
}

func TestRoomStorage_ListRooms_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	storage := NewRoomStorage(db)

	mock.ExpectQuery("SELECT id, name, description, capacity, created_at FROM rooms").
		WillReturnError(sql.ErrConnDone)

	_, err = storage.RoomList()
	assert.Error(t, err)
}

func TestScheduleStorage_CreateSchedule_RoomNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	storage := NewScheduleStorage(db)
	roomID := "non-existent-id"

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(roomID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	_, err = storage.CreateSchedule(roomID, []int{1}, "09:00", "18:00")
	assert.Error(t, err)
	assert.Equal(t, "room does not exist", err.Error())
}

func TestScheduleStorage_CreateSchedule_AlreadyExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	storage := NewScheduleStorage(db)
	roomID := "test-room-id"

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(roomID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(roomID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	_, err = storage.CreateSchedule(roomID, []int{1}, "09:00", "18:00")
	assert.Error(t, err)
	assert.Equal(t, "schedule already exists", err.Error())
}

func TestRoomStorage_CreateRoom_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	storage := NewRoomStorage(db)

	mock.ExpectExec("INSERT INTO rooms").
		WillReturnError(sql.ErrConnDone)

	_, err = storage.CreateRoom("Test", nil, nil)
	assert.Error(t, err)
}

func TestRoomStorage_ListRooms_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	storage := NewRoomStorage(db)

	mock.ExpectQuery("SELECT id, name, description, capacity, created_at FROM rooms").
		WillReturnError(sql.ErrConnDone)

	_, err = storage.RoomList()
	assert.Error(t, err)
}

func TestParseTime_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		timeStr  string
		wantHour int
		wantMin  int
	}{
		{"single digit hour", "9:05", 9, 5},
		{"leading zero", "09:05", 9, 5},
		{"midnight", "00:00", 0, 0},
		{"end of day", "23:59", 23, 59},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hour, min := parseTime(tt.timeStr)
			assert.Equal(t, tt.wantHour, hour)
			assert.Equal(t, tt.wantMin, min)
		})
	}
}
