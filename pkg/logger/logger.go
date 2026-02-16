package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	logger zerolog.Logger
}

func New(level string, output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        output,
		TimeFormat: time.RFC3339,
	}

	logLevel := zerolog.InfoLevel
	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	}

	logger := zerolog.New(consoleWriter).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{logger: logger}
}

func Default() *Logger {
	return New("info", os.Stdout)
}

func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	event := l.logger.Debug()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func (l *Logger) Info(msg string, fields map[string]interface{}) {
	event := l.logger.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	event := l.logger.Warn()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func (l *Logger) Error(msg string, err error, fields map[string]interface{}) {
	event := l.logger.Error()
	if err != nil {
		event = event.Err(err)
	}
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func (l *Logger) Fatal(msg string, err error, fields map[string]interface{}) {
	event := l.logger.Fatal()
	if err != nil {
		event = event.Err(err)
	}
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func (l *Logger) With(fields map[string]interface{}) *Logger {
	logger := l.logger.With()
	for k, v := range fields {
		logger = logger.Interface(k, v)
	}
	return &Logger{logger: logger.Logger()}
}

func SetGlobalLogger(l *Logger) {
	log.Logger = l.logger
}
