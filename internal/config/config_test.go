package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	origDBHost := os.Getenv("DB_HOST")
	origServerPort := os.Getenv("SERVER_PORT")
	defer func() {
		os.Setenv("DB_HOST", origDBHost)
		os.Setenv("SERVER_PORT", origServerPort)
	}()

	os.Setenv("DB_HOST", "testhost")
	os.Setenv("SERVER_PORT", "9090")

	cfg := Load()
	assert.Equal(t, "testhost", cfg.DBHost)
	assert.Equal(t, "9090", cfg.ServerPort)

	assert.Equal(t, "secret", cfg.JWTSecret)
	assert.Equal(t, 24, cfg.JWTTTL)
}

func TestLoadDefaults(t *testing.T) {
	origDBHost := os.Getenv("DB_HOST")
	origServerPort := os.Getenv("SERVER_PORT")
	origJWTSecret := os.Getenv("JWT_SECRET")
	defer func() {
		os.Setenv("DB_HOST", origDBHost)
		os.Setenv("SERVER_PORT", origServerPort)
		os.Setenv("JWT_SECRET", origJWTSecret)
	}()

	os.Unsetenv("DB_HOST")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("JWT_SECRET")

	cfg := Load()
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "8080", cfg.ServerPort)
	assert.Equal(t, "secret", cfg.JWTSecret)
	assert.Equal(t, 24, cfg.JWTTTL)
}

func TestGetEnv(t *testing.T) {
	orig := os.Getenv("TEST_KEY")
	defer os.Setenv("TEST_KEY", orig)

	os.Setenv("TEST_KEY", "test_value")
	result := getEnv("TEST_KEY", "default")
	assert.Equal(t, "test_value", result)

	os.Unsetenv("TEST_KEY")
	result = getEnv("TEST_KEY", "default")
	assert.Equal(t, "default", result)
}

func TestGetEnvAsInt(t *testing.T) {
	orig := os.Getenv("TEST_INT")
	defer os.Setenv("TEST_INT", orig)

	os.Setenv("TEST_INT", "123")
	result := getEnvInt("TEST_INT", 0)
	assert.Equal(t, 123, result)

	os.Unsetenv("TEST_INT")
	result = getEnvInt("TEST_INT", 42)
	assert.Equal(t, 42, result)

	os.Setenv("TEST_INT", "notanumber")
	result = getEnvInt("TEST_INT", 10)
	assert.Equal(t, 10, result)
}
