package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	logLevel := logrus.InfoLevel
	if logLevelRaw := os.Getenv("LOG_LEVEL"); logLevelRaw != "" {
		level, err := logrus.ParseLevel(logLevelRaw)
		if err != nil {
			panic("can not set up logger, invalid LOG_LEVEL")
		}
		logLevel = level
	}

	return &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			DisableTimestamp: true,
		},
		Level:        logLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
}
