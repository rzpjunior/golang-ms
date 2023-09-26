package order

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockApplicant struct {
	mock.Mock
}

func (m *MockApplicant) Messages() error {
	args := m.Called()
	return args.Error(0)

}
func TestMessage(t *testing.T) {
	a := new(MockApplicant)
	a.On("Messages").Return(nil)

}
func TestSetFinishAt(t *testing.T) {
	input := 2
	expected := time.Now()
	expectedDate := expected.Format("2006-01-02")
	var finishedAt string
	if input == 2 {
		finishedAt = time.Now().Format("2006-01-02")
	}
	assert.Equal(t, expectedDate, finishedAt)
}
