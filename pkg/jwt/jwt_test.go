package jwt

import (
	"testing"
)

func TestGenerateToken(t *testing.T) {
	secret := "test-secret"
	userID := "test-user"
	role := "user"
	ttl := 24

	token, err := GenerateToken(userID, role, secret, ttl)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}
}

func TestParseToken(t *testing.T) {
	secret := "test-secret"
	userID := "test-user"
	role := "admin"
	ttl := 24

	token, _ := GenerateToken(userID, role, secret, ttl)

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if claims.UserId != userID {
		t.Errorf("Expected userID %s, got %s", userID, claims.UserId)
	}
	if claims.Role != role {
		t.Errorf("Expected role %s, got %s", role, claims.Role)
	}
}

func TestParseTokenInvalid(t *testing.T) {
	secret := "test-secret"
	invalidToken := "invalid.token.string"

	_, err := ParseToken(invalidToken, secret)
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	secret := "test-secret"
	wrongSecret := "wrong-secret"
	userID := "test-user"
	role := "user"
	ttl := 24

	token, _ := GenerateToken(userID, role, secret, ttl)

	_, err := ParseToken(token, wrongSecret)
	if err == nil {
		t.Error("Expected error for wrong secret")
	}
}

func TestParseTokenExpired(t *testing.T) {
	secret := "test-secret"
	userID := "test-user"
	role := "user"
	ttl := -1

	token, _ := GenerateToken(userID, role, secret, ttl)

	_, err := ParseToken(token, secret)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}
