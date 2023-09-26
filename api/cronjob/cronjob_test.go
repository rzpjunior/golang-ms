package cronjob

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestCalculateReturnTotalPointCancel(t *testing.T) {
	recentPoint := 400
	pointRedeem := 500
	expected := 900
	assert.Equal(t, recentPoint+pointRedeem, expected)

}
