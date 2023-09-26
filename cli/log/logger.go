package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = New()

func New(prefix ...string) *logrus.Logger {
	l := logrus.Logger{
		Out:       os.Stderr,
		Formatter: NewFormater(true, prefix...),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

	return &l
}
