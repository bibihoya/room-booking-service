package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgresDB_InvalidConfig(t *testing.T) {
	_, err := NewPostgresDB("", "", "", "", "", "")
	assert.Error(t, err)
}
