// Package logger fornece logging estruturado para toda a aplicação.
package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger é um wrapper em torno do zerolog.Logger
type Logger struct {
	logger zerolog.Logger
}

// New cria uma nova instância de Logger
func New(level string, output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}

	// Configurar output formatado para desenvolvimento
	consoleWriter := zerolog.ConsoleWriter{
		Out:        output,
		TimeFormat: time.RFC3339,
	}

	// Determinar nível de log
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

// Default retorna o logger padrão da aplicação
func Default() *Logger {
	return New("info", os.Stdout)
}

// Debug registra uma mensagem de debug
func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	event := l.logger.Debug()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Info registra uma mensagem informativa
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	event := l.logger.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Warn registra uma mensagem de aviso
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	event := l.logger.Warn()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Error registra uma mensagem de erro
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

// Fatal registra uma mensagem fatal e encerra a aplicação
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

// With cria um logger filho com campos adicionais
func (l *Logger) With(fields map[string]interface{}) *Logger {
	logger := l.logger.With()
	for k, v := range fields {
		logger = logger.Interface(k, v)
	}
	return &Logger{logger: logger.Logger()}
}

// SetGlobalLogger define o logger global do zerolog
func SetGlobalLogger(l *Logger) {
	log.Logger = l.logger
}
