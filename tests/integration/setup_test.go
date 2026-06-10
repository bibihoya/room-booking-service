package integration

import (
	"os"
	"testing"

	"github.com/internships-backend/test-backend-bibihoya-1/tests"
)

func TestMain(m *testing.M) {
	tests.InitTestDB()

	code := m.Run()

	tests.CloseTestDB()

	os.Exit(code)
}
