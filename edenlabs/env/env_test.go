package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetConfigSuccess(t *testing.T) {
	cfg, err := Env("env")
	host := cfg.GetString("app.host")
	assert.Equal(t, host, "0.0.0.0")
	assert.NotNil(t, cfg)
	assert.Nil(t, err)
}

func Test_GetConfigFailed(t *testing.T) {
	cfg, err := Env("unknown")
	host := cfg.GetString("app.host")
	assert.Equal(t, host, "")
	assert.NotNil(t, err)
}
