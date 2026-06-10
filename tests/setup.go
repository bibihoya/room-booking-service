package tests

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/internships-backend/test-backend-bibihoya-1/internal/database"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/storage"
	_ "github.com/lib/pq"
)

var (
	TestDB          *sql.DB
	UserStorage     *storage.UserStorage
	RoomStorage     *storage.RoomStorage
	ScheduleStorage *storage.ScheduleStorage
	SlotStorage     *storage.SlotStorage
	BookingStorage  *storage.BookingStorage
)

func InitTestDB() {
	time.Sleep(2 * time.Second)

	connStr := "host=localhost port=5432 user=postgres password=postgres sslmode=disable"
	adminDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer adminDB.Close()

	if err := adminDB.Ping(); err != nil {
		log.Fatalf("Cannot ping postgres: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	dbName := "meeting_booking_test"

	_, err = adminDB.Exec(fmt.Sprintf(`
		SELECT pg_terminate_backend(pid) 
		FROM pg_stat_activity 
		WHERE datname = '%s' AND pid <> pg_backend_pid()
	`, dbName))
	if err != nil {
		log.Printf("Warning: failed to terminate connections: %v", err)
	}

	_, err = adminDB.Exec("DROP DATABASE IF EXISTS " + dbName)
	if err != nil {
		log.Printf("Warning: failed to drop test DB: %v", err)
	}

	_, err = adminDB.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		log.Fatalf("Failed to create test DB: %v", err)
	}
	log.Printf("Test database created: %s", dbName)

	TestDB, err = database.NewPostgresDB(
		"localhost", "5432", "postgres", "postgres", dbName, "disable",
	)
	if err != nil {
		log.Fatalf("Failed to connect to test DB: %v", err)
	}

	if err := database.Migrate(TestDB); err != nil {
		log.Fatalf("Failed to migrate test DB: %v", err)
	}
	log.Println("Migrations completed")

	if err := insertTestUsers(); err != nil {
		log.Fatalf("Failed to insert test users: %v", err)
	}
	log.Println("Test users inserted")

	RoomStorage = storage.NewRoomStorage(TestDB)
	ScheduleStorage = storage.NewScheduleStorage(TestDB)
	SlotStorage = storage.NewSlotStorage(TestDB)
	BookingStorage = storage.NewBookingStorage(TestDB)
}

func insertTestUsers() error {
	_, err := TestDB.Exec(`
		DELETE FROM users WHERE id IN (
			'11111111-1111-1111-1111-111111111111',
			'22222222-2222-2222-2222-222222222222'
		)
	`)
	if err != nil {
		return err
	}

	_, err = TestDB.Exec(`
		INSERT INTO users (id, email, role) VALUES 
			('11111111-1111-1111-1111-111111111111', 'admin@example.com', 'admin'),
			('22222222-2222-2222-2222-222222222222', 'user@example.com', 'user')
	`)
	return err
}

func CleanDB() {
	if TestDB == nil {
		return
	}
	tables := []string{"bookings", "slots", "schedules", "rooms"}
	for _, table := range tables {
		_, err := TestDB.Exec("TRUNCATE TABLE " + table + " CASCADE")
		if err != nil {
			log.Printf("Warning: failed to truncate table %s: %v", table, err)
		}
	}
}

func CloseTestDB() {
	if TestDB != nil {
		TestDB.Close()
	}
}
