package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMysqlDatabase(t *testing.T) {
	opt := DBMysqlOption{}
	err := NewMysqlDatabase("", nil, "", opt)
	assert.NotNil(t, err)
}
