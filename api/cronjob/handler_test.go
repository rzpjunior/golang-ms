package cronjob

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestRounding(t *testing.T) {
	// assert equality
	case1 := 1234005.00
	sumCalculationPoints := math.Round(math.Floor(case1 / 100))

	assert.Equal(t, float64(12340), sumCalculationPoints, "they should be equal for rounding")
}
