// internal/pkg/logger/zap.go
package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strings"
)

type LogrusLogger struct {
	logger *logrus.Logger
}

func NewLogrus() Logger {
	l := logrus.New()

	l.SetLevel(logrus.InfoLevel)

	// Jangan pakai ini kalau logger kamu dibungkus wrapper.
	// Karena akan menunjuk ke LogrusLogger.Info().
	l.SetReportCaller(false)

	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05-07:00",
	})

	return &LogrusLogger{
		logger: l,
	}
}

func (l *LogrusLogger) Info(args ...interface{}) {
	l.logger.WithFields(getCallerFields(2)).Info(args...)
}

func (l *LogrusLogger) Error(args ...interface{}) {
	l.logger.WithFields(getCallerFields(2)).Error(args...)
}

func (l *LogrusLogger) Warn(args ...interface{}) {
	l.logger.WithFields(getCallerFields(2)).Warn(args...)
}

func (l *LogrusLogger) Debug(args ...interface{}) {
	l.logger.WithFields(getCallerFields(2)).Debug(args...)
}

func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.logger.WithFields(getCallerFields(2)).Infof(format, args...)
}

func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.WithFields(getCallerFields(2)).Errorf(format, args...)
}

func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.WithFields(getCallerFields(2)).Warnf(format, args...)
}

func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.WithFields(getCallerFields(2)).Debugf(format, args...)
}

func getCallerFields(skip int) logrus.Fields {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return logrus.Fields{
			"file": "unknown",
			"func": "unknown",
		}
	}

	fn := runtime.FuncForPC(pc)

	funcName := "unknown"
	if fn != nil {
		funcName = trimFunctionName(fn.Name())
	}

	return logrus.Fields{
		"file": fmt.Sprintf("%s:%d", filepath.Base(file), line),
		"func": funcName,
	}
}

func trimFunctionName(name string) string {
	// Optional: biar tidak terlalu panjang.
	// Dari:
	// saetechnology-be/internal/delivery/http/middleware.(*LoggerMiddleware).ServeHTTP
	// Jadi:
	// middleware.(*LoggerMiddleware).ServeHTTP

	parts := strings.Split(name, "/")
	if len(parts) == 0 {
		return name
	}

	return parts[len(parts)-1]
}
