package loggertest

import (
	"testing"

	"github.com/jeffalyanak/check_freenas_api/logger"
)

func TestAbs(t *testing.T) {
	_, err := logger.Get()

	if err != nil {
		t.Errorf("Abs(-1) = %d; want 1", err)
	}
}
