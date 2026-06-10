package database

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255),
			role VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS rooms (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			capacity INT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS schedules (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			room_id UUID REFERENCES rooms(id) ON DELETE CASCADE,
			days_of_week INTEGER[] NOT NULL,
			start_time VARCHAR(5) NOT NULL,
			end_time VARCHAR(5) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS slots (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			room_id UUID REFERENCES rooms(id),
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS bookings (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			slot_id UUID REFERENCES slots(id),
			user_id UUID REFERENCES users(id),
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			conference_link VARCHAR(512),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		`CREATE INDEX IF NOT EXISTS idx_slots_room_time ON slots(room_id, start_time)`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_user_status ON bookings(user_id, status)`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_slot_status ON bookings(slot_id, status) WHERE status = 'active'`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	insertUsers := `INSERT INTO users (id, email, role) VALUES 
		('11111111-1111-1111-1111-111111111111', 'admin@example.com', 'admin'),
		('22222222-2222-2222-2222-222222222222', 'user@example.com', 'user')
		ON CONFLICT (id) DO NOTHING`

	if _, err := db.Exec(insertUsers); err != nil {
		return fmt.Errorf("failed to insert test users: %w", err)
	}

	return nil
}
