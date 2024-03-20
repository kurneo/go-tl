package log

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"sync"
)

var (
	instance  Contract
	once      sync.Once
	file      *os.File
	formatter = &logrus.JSONFormatter{}
)

type Contract interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

func New() (Contract, error) {
	var err error
	c := viper.GetString("LOG_DEFAULT_CHANNEL")
	if c == "" {
		c = "daily"
	}
	once.Do(func() {
		l := newLogrus(logrus.InfoLevel)
		switch c {
		case "daily":
			instance, err = newDailyDriver(l)
			break
		case "single":
			instance, err = newSingleDriver(l)
			break
		case "stdout":
			instance = newStdoutDriver(l)
			break
		default:
			err = errors.New("log channel is invalid")
		}
		setupHook(l)
	})
	return instance, err
}
