package integration

import (
	"testing"

	"github.com/internships-backend/test-backend-bibihoya-1/tests"
)

func TestRegisterSuccess(t *testing.T) {
	tests.CleanDB()

	user, err := tests.UserStorage.CreateUser("test@example.com", "123456", "user")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == "" {
		t.Error("User ID should not be empty")
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.Email)
	}
	if user.Role != "user" {
		t.Errorf("Expected role user, got %s", user.Role)
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	tests.CleanDB()

	_, err := tests.UserStorage.CreateUser("test@example.com", "123456", "user")
	if err != nil {
		t.Fatalf("First creation failed: %v", err)
	}

	_, err = tests.UserStorage.CreateUser("test@example.com", "123456", "user")
	if err == nil {
		t.Error("Expected error for duplicate email")
	}
}

func TestLoginSuccess(t *testing.T) {
	tests.CleanDB()

	_, err := tests.UserStorage.CreateUser("test@example.com", "123456", "user")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user, err := tests.UserStorage.GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.Email)
	}

	if !tests.UserStorage.CheckPassword(user, "123456") {
		t.Error("Password check failed")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	tests.CleanDB()

	_, err := tests.UserStorage.CreateUser("test@example.com", "123456", "user")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user, err := tests.UserStorage.GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if tests.UserStorage.CheckPassword(user, "wrongpassword") {
		t.Error("Password check should fail for wrong password")
	}
}

func TestLoginUserNotFound(t *testing.T) {
	tests.CleanDB()

	_, err := tests.UserStorage.GetUserByEmail("nonexistent@example.com")
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
}
