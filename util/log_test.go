package util

import (
	"testing"
	"time"
)

func TestLogWithRed(t *testing.T) {
	LogWithRed("testLog" + time.Now().String())
}
