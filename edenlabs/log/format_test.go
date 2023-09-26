package log

import (
	"bytes"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormater(t *testing.T) {
	var output bytes.Buffer

	f := LogFormatter{
		Tty:         false,
		Prefix:      "LOGGER",
		ForceColors: false,
	}

	log := &logrus.Logger{
		Out:       &output,
		Formatter: &f,
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	log.WithFields(logrus.Fields{
		"prefix": "OK",
		"file":   1,
		"time":   1,
		"msg":    1,
		"level":  1,
		"logger": 1,
		"error":  errors.New("error here"),
	}).Error("OK")

	assert.Contains(t, output.String(), "fields.time")
	assert.Contains(t, output.String(), "fields.msg")
	assert.Contains(t, output.String(), "fields.level")
	assert.Contains(t, output.String(), "fields.logger")
}
