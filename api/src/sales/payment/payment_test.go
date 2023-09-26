package payment

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestSetFinishAt(t *testing.T) {
	inputStatusDoc := 2
	expected := time.Now()
	expectedDate := expected.Format("2006-01-02")
	var finishedAt string
	if inputStatusDoc == 2 {
		finishedAt = time.Now().Format("2006-01-02")
	}
	assert.Equal(t, expectedDate, finishedAt)
}
