package generator

import (
	"math/rand"
	"time"
)

func DummyTime() (data time.Time) {
	data, _ = time.Parse("2006-01-02", "2023-01-01")
	return
}

func DummyInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func DummyInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

func DummyFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
