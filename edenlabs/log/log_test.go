package log

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLogger(t *testing.T) {
	lgr := NewLogger("any-service")

	assert.NotNil(t, lgr)
}

func TestAddMessage(t *testing.T) {
	lgr := NewLogger("any-service")
	lgr.AddMessage(logs.LevelError, "error level")
	lgr.AddMessage(logs.LevelWarning, "error level")
	lgr.AddMessage(logs.LevelDebug, "error level")
	lgr.AddMessage(logs.LevelInfo, "error level")

	assert.NotNil(t, lgr)
}

func TestNewChildLogger(t *testing.T) {
	lgr := NewLogger("any-service")
	lgrchild := lgr.NewChildLogger()

	assert.NotNil(t, lgrchild)
}

func TestLogger_Print(t *testing.T) {
	lgr := NewLogger("any-service")
	lgr.AddMessage(logs.LevelError, "error level").Print("Error")

	assert.NotNil(t, lgr)
}
