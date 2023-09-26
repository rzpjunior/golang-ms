package log

import (
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type Level uint32

const (
	FieldLogID           = "log_id"
	FieldEndpoint        = "endpoint"
	FieldMethod          = "method"
	FieldServiceName     = "service"
	FieldRequestBody     = "request_body"
	FieldRequestHeaders  = "request_headers"
	FieldResponseBody    = "response_body"
	FieldResponseHeaders = "response_headers"
)

type message struct {
	Message  interface{} `json:"message"`
	Level    Level       `json:"level"`
	File     string      `json:"file"`
	FuncName string      `json:"func"`
	Line     int         `json:"line"`
}

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

func (level Level) MarshalText() ([]byte, error) {
	switch level {
	case TraceLevel:
		return []byte("trace"), nil
	case DebugLevel:
		return []byte("debug"), nil
	case InfoLevel:
		return []byte("info"), nil
	case WarnLevel:
		return []byte("warning"), nil
	case ErrorLevel:
		return []byte("error"), nil
	case FatalLevel:
		return []byte("fatal"), nil
	case PanicLevel:
		return []byte("panic"), nil
	}

	return nil, fmt.Errorf("not a valid logrus level %d", level)
}

type Logger struct {
	logger  *log.Logger
	fields  sync.Map
	id      string
	service string
}

func (l *Logger) NewChildLogger() (logger *Logger) {
	logger = newLogger(l.service, l.id, nil)
	return
}

func (l *Logger) AddMessage(level Level, message ...interface{}) *Logger {
	l.setCaller(level, 2, message...)
	return l
}

func (l *Logger) Print(directMsg ...interface{}) {
	if len(directMsg) > 0 {
		l.setCaller(DebugLevel, 2, directMsg...)
	}

	stackVal, _ := l.fields.Load("stack")
	messages := ensureStackType(stackVal)

	if len(messages) > 0 {
		maxLevel := l.findMaxLevel(messages)
		if maxLevel < WarnLevel {
			l.logger.SetOutput(os.Stderr)
		} else {
			l.logger.SetOutput(os.Stdout)
		}

		entry := l.logger.WithFields(l.syncMapToLogFields())
		entry.Logf(log.Level(maxLevel), "%+v", messages[0].Message)
	}

	l.clear()
}

func (l *Logger) Logger() (log *log.Logger) {
	log = l.logger
	return
}

func (l *Logger) findMaxLevel(msgs []message) (maxLevel Level) {
	currentMaxLevel := TraceLevel
	for _, msg := range msgs {
		currentMaxLevel = Level(math.Min(float64(msg.Level), float64(currentMaxLevel)))
	}

	maxLevel = currentMaxLevel
	return
}

func (l *Logger) clear() {
	l.fields.Range(func(key interface{}, value interface{}) bool {
		k, ok := key.(string)
		if ok {
			if k != FieldLogID {
				l.fields.Delete(key)
			}

			return true
		}

		return false
	})

	id := uuid.NewV1().String()

	l.fields.Store(FieldLogID, id)
	l.fields.Store("stack", []message{})
}

func (l *Logger) syncMapToLogFields() (fields log.Fields) {
	fields = make(log.Fields)

	l.fields.Range(func(key interface{}, value interface{}) bool {
		k, ok := key.(string)
		if ok {
			fields[k] = value
			return true
		}

		return false
	})

	return
}

func (l *Logger) setCaller(level Level, callerLevel int, msgs ...interface{}) {
	x := len(msgs)
	if msgs == nil || x == 0 {
		return
	}

	for _, val := range msgs {
		if val == "" {
			continue
		}

		if pc, file, line, ok := runtime.Caller(callerLevel); ok {
			fName := runtime.FuncForPC(pc).Name()

			err, ok := val.(error)
			if ok && err != nil {
				val = err.Error()
			}

			vmsg := message{
				Message:  val,
				Level:    level,
				File:     file,
				FuncName: fName,
				Line:     line,
			}

			l.addMessageStack(vmsg)
		}
	}
}

func (l *Logger) addMessageStack(msg ...message) {
	stackVal, _ := l.fields.Load("stack")
	stack := ensureStackType(stackVal)
	stack = append(stack, msg...)

	l.fields.Store("stack", stack)
}

func ensureStackType(stack interface{}) (val []message) {
	val, ok := stack.([]message)
	if !ok {
		return []message{}
	}

	return
}

func newLog(formatter log.Formatter, out io.Writer, level log.Level, reportCaller bool) (l *log.Logger) {
	l = log.New()
	l.SetFormatter(formatter)
	l.SetOutput(out)
	l.SetLevel(level)

	return
}

func newLogger(serviceName string, logID string, newLogger *log.Logger) (logger *Logger) {
	if newLogger == nil {
		formatter := &log.JSONFormatter{
			TimestampFormat: time.RFC3339,
			// PrettyPrint:     true,
			FieldMap: log.FieldMap{
				log.FieldKeyMsg: "log_message",
			},
		}

		newLogger = newLog(formatter, os.Stdout, log.TraceLevel, false)
	}

	logger = new(Logger)
	logger.logger = newLogger
	logger.fields.Range(func(key interface{}, value interface{}) bool {
		logger.fields.Delete(key)
		return true
	})
	logger.fields.Store("stack", []message{})

	id := uuid.NewV1().String()
	if logID != "" {
		id = logID
	}

	logger.fields.Store(FieldLogID, id)
	logger.fields.Store(FieldServiceName, serviceName)
	logger.id = id
	logger.service = serviceName
	return
}

func NewLogger(serviceName string) (logger *Logger) {
	logger = newLogger(serviceName, "", nil)
	return
}

func NewLoggerWithClient(serviceName string, client *log.Logger) (logger *Logger) {
	logger = newLogger(serviceName, "", client)
	return
}
