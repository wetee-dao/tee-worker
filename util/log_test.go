package util

import (
	"testing"
	"time"
)

func TestLogError(t *testing.T) {
	LogError("testLog" + time.Now().String())
}
