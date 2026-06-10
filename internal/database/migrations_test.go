package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMigrate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS rooms").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schedules").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS slots").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS bookings").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE INDEX IF NOT EXISTS idx_slots_room_time").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE INDEX IF NOT EXISTS idx_bookings_user_status").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE INDEX IF NOT EXISTS idx_bookings_slot_status").WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(0, 0))

	err = Migrate(db)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
