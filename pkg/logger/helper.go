package logger

import (
	"fmt"
	"github.com/kurneo/go-template/config"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

func newLogrus(level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(formatter)
	logger.SetLevel(level)
	return logger
}

func createLogFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
}

func getLogLevel(l string) logrus.Level {
	var level logrus.Level

	switch strings.ToLower(l) {
	case "error":
		level = logrus.ErrorLevel
	case "warn":
		level = logrus.WarnLevel
	case "info":
		level = logrus.InfoLevel
	case "debug":
		level = logrus.DebugLevel
	default:
		level = logrus.InfoLevel
	}

	return level
}

func getCalledFile(skip int) string {
	_, callerFile, line, ok := runtime.Caller(skip)
	if !ok {
		callerFile = "<???>"
		line = 1
	}
	return fmt.Sprintf("%s:%d", callerFile, line)
}

func setupHook(cfg config.Log, l *logrus.Logger) {
	//Telegram
	if c, ok := cfg.Hooks["telegram"]; ok {
		if c.Get("enable").(bool) {
			addTelegramHook(c, l)
		}
	}
}
