package timex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToStartTime(t *testing.T) {
	dateStr := "2023-01-01"

	var baseDate time.Time
	var err error

	baseDate, err = time.Parse(InFormatDate, dateStr)
	assert.Nil(t, err)
	dateFrom := ToStartTime(baseDate)
	assert.Equal(t, dateFrom.Format(InFormatDateTime), "2023-01-01 00:00:00")

	baseDate, err = time.Parse(InFormatDate, dateStr)
	assert.Nil(t, err)
	assert.Nil(t, err)
	dateTo := ToLastTime(baseDate)
	assert.Equal(t, dateTo.Format(InFormatDateTime), "2023-01-01 23:59:59")
}
