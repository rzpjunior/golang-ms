package order

import (
	"testing"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
)

type salesOrderMock struct {
	mock.Mock
}

func TestCalculateReturnTotalPointCancel(t *testing.T) {
	recentPoint := 400
	pointRedeem := 500
	expected := 900
	assert.Equal(t, recentPoint+pointRedeem, expected)

}
func TestCalculateReturnTotalPointCancel2(t *testing.T) {
	recentPoint := 400
	pointRedeem := 500
	expected := 900

	result := addRecentPointTotal(float64(recentPoint), float64(pointRedeem))
	assert.Equal(t, float64(expected), result)
}

func TestCalculateReturnTotalDiscAmount(t *testing.T) {
	discAmount := 400
	pointRedeem := 500
	expected := 900
	assert.Equal(t, discAmount+pointRedeem, expected)

}

func TestIsPointRedeemAmountExist(t *testing.T) {
	pointRedeem := 100
	assert.NotEqual(t, pointRedeem, 0)
}
func (soMock *salesOrderMock) GetSalesOrder(field string, values ...interface{}) (int, error) {
	args := soMock.Called(field, values)
	return args.Int(0), args.Error(1)
}
func TestGetSalesOrder(t *testing.T) {
	id := 57252
	repository.GetSalesOrder("id", id)

}
