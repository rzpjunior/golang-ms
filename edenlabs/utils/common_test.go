package utils

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncrypt(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    interface{}
		expected string
	}{
		{1, "65536"},
		{4040, "264765440"},
		{"264765440", "4040"},
		{"65536", "1"},
		{"randomstring", "0"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, Encrypt(test.param))
	}
}

func TestDecrypt(t *testing.T) {
	var tests = []struct {
		param    interface{}
		expected int64
		valid    bool
	}{
		{"264765440", 4040, true},
		{"65536", 1, true},
		{"randomstring", 0, false},
	}

	for _, test := range tests {
		v, e := Decrypt(test.param)
		if e != nil {
			_, ok := e.(*DecryptionError)
			if !ok {
				assert.Fail(t, "Error type is invalid")
			}
		}

		assert.Equal(t, test.expected, v)
		assert.Equal(t, e == nil, test.valid)
	}
}

func TestRandomStr(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		len int
	}{
		{10},
		{4},
		{1000},
		{5},
	}

	for _, test := range tests {
		str := RandomStr(test.len)
		assert.Equal(t, test.len, len(str), "Not equal length given.")
	}
}

func TestRandomNumeric(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		len int
	}{
		{10},
		{4},
		{1000},
		{5},
	}

	for _, test := range tests {
		str := RandomNumeric(test.len)
		assert.Equal(t, test.len, len(str), "Not equal length given.")
	}
}

func TestPasswordHash(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param string
	}{
		{"randompassword"},
		{"123456"},
	}
	for _, test := range tests {
		hash, e := PasswordHasher(test.param)
		assert.NoError(t, e)

		match := PasswordHash(hash, test.param)
		assert.NoError(t, match)
	}

}

func TestContains(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		slicer   []string
		value    string
		expected bool
	}{
		{[]string{"one", "two"}, "one", true},
		{[]string{"one", "two"}, "three", false},
	}

	for _, test := range tests {
		res := Contains(test.slicer, test.value)
		assert.Equal(t, test.expected, res)
	}

	strArr := []string{"a", "b", "c"}
	res := Contains(strArr, "a")
	assert.Equal(t, res, true)
	res2 := Contains(strArr, "x")
	assert.Equal(t, res2, false)
}

func TestFloatPrecision(t *testing.T) {
	d := 10.123123123
	r := FloatPrecision(d, 2)
	assert.Equal(t, 10.12, r)
}

func TestRounder(t *testing.T) {
	d := Rounder(2.588, 0.5, 1)
	assert.Equal(t, float64(2.6), d)

	d2 := Rounder(2.588, 0.5, 0)
	assert.Equal(t, float64(3), d2)

	d3 := Rounder(2.588, 0.9, 1)
	assert.Equal(t, float64(2.5), d3)

	d4 := Rounder(-38.288888, 0.5, 2)
	assert.Equal(t, float64(-38.28), d4)

	d5 := Rounder(3.333333, 0.3, 1)
	assert.Equal(t, float64(3.4), d5)

	d6 := Rounder(3.1, 0.1, 2)
	assert.Equal(t, float64(3.1), d6)
}

type ModelTest struct {
	Image         string `orm:"column(image);null" json:"image"`
	BarcodeType   string `orm:"column(barcode_type);null;options(qr_code,ean_13,ean_8,upc_a,upc_e)" json:"barcode_type"`
	BarcodeNumber string `orm:"column(barcode_number);size(50);null" json:"barcode_number"`
	BarcodeImage  string `orm:"column(barcode_image);null" json:"barcode_image"`
	Note          string `orm:"column(note);null" json:"note"`
	Test          string `orm:"-" json:"attributes"`
}

func TestFields(t *testing.T) {
	x := Fields(ModelTest{}, "note")
	assert.Equal(t, 4, len(x))

	i := &ModelTest{
		Image: "test",
	}

	x = Fields(i, "note")
	assert.Equal(t, 4, len(x))
}

func TestDecryptionError_Error(t *testing.T) {
	descErrorNil := &DecryptionError{}

	assert.Equal(t, descErrorNil, &DecryptionError{})

	descError := &DecryptionError{
		Message: "test",
		Values:  "val",
	}
	assert.NotNil(t, descError)

	val := descError.Error()

	assert.Equal(t, val, "Invalid encryption values: val")

}

func TestPasswordHasher(t *testing.T) {
	var err error
	err = PasswordHash("xxx", "pass")

	assert.Equal(t, err, errors.New("crypto/bcrypt: hashedSecret too short to be a bcrypted password"))

	err = PasswordHash("xxx-log-hasher", "pass1234")
	assert.NotNil(t, err)

	h, err := PasswordHasher("pass1234")
	fmt.Println(h)
	assert.Nil(t, err)
	assert.NotNil(t, h)
}
