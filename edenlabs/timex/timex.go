package timex

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/constants"
)

const (
	InFormatDateTime = "2006-01-02 15:04:05"
	InFormatDate     = "2006-01-02"
	InFormatTime     = "15:04:05"
)

func ToLocTime(ctx context.Context, baseTime time.Time) (locTime time.Time) {
	valueKeyTimezone := ctx.Value(constants.KeyTimezone)
	if valueKeyTimezone != nil {
		timezone := ctx.Value(constants.KeyTimezone).(string)
		loc, _ := time.LoadLocation(timezone)
		locTime = baseTime.In(loc)
	} else {
		loc, _ := time.LoadLocation("")
		locTime = baseTime.In(loc)
	}
	return
}

func ToStartTime(baseDate time.Time) time.Time {
	return baseDate.Add(time.Hour*time.Duration(0) +
		time.Minute*time.Duration(0) +
		time.Second*time.Duration(0))
}

func ToLastTime(baseDate time.Time) time.Time {
	return baseDate.Add(time.Hour*time.Duration(23) +
		time.Minute*time.Duration(59) +
		time.Second*time.Duration(59))
}

func IsValid(validTime time.Time) bool {
	if validTime.Year() == 1 || validTime.Year() == 1970 {
		return false
	}
	return true
}
