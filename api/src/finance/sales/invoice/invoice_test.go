package invoice

import (
	"testing"

	"github.com/go-playground/assert"
)

func TestIsPointRedeemAmountExist(t *testing.T) {
	pointRedeem := 100
	assert.NotEqual(t, pointRedeem, 0)
}
func TestCalculateReturnTotalDiscAmount(t *testing.T) {
	discAmount := 400
	pointRedeem := 500
	expected := 900
	assert.Equal(t, discAmount+pointRedeem, expected)

}
