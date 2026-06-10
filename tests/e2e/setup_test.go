package e2e

import (
	"os"
	"testing"

	"github.com/internships-backend/test-backend-bibihoya-1/tests"
)

var initialized = false

func TestMain(m *testing.M) {
	if !initialized {
		tests.InitTestDB()
		initialized = true
	}

	code := m.Run()

	os.Exit(code)
}
