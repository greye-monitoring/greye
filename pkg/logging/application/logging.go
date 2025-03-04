package application

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"greye/pkg/logging/domain/ports"
	"os"
)

type Logger struct {
	logrus *logrus.Logger
}

var _ ports.LoggerApplication = (*Logger)(nil)

func NewLogger(logLevel string) *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(level)
	}
	return &Logger{logrus: log}
}
func (l Logger) Trace(msg string, args ...interface{}) {
	if len(args) > 0 {
		message := fmt.Sprintf(msg, args...)
		l.logrus.Trace(message)
	} else {
		l.logrus.Trace(msg)
	}
}

func (l Logger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		message := fmt.Sprintf(msg, args...)
		l.logrus.Debug(message)
	} else {
		l.logrus.Debug(msg)
	}
}

func (l Logger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		message := fmt.Sprintf(msg, args...)
		l.logrus.Info(message)
	} else {
		l.logrus.Info(msg)
	}
}

func (l Logger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		message := fmt.Sprintf(msg, args...)
		l.logrus.Warn(message)
	} else {
		l.logrus.Warn(msg)
	}
}

func (l Logger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		message := fmt.Sprintf(msg, args...)
		l.logrus.Error(message)
	} else {
		l.logrus.Error(msg)
	}
}
