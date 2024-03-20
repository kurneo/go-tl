package log

import (
	"errors"
	"github.com/sirupsen/logrus"
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

type Config struct {
	Channel        string
	Daily          DailyConfig
	Singe          SingeConfig
	StdOut         StdOutConfig
	TeleHookConfig TeleHookConfig
}

func New(c Config) (Contract, error) {
	var err error
	once.Do(func() {
		l := newLogrus(logrus.InfoLevel)
		switch c.Channel {
		case "daily":
			instance, err = newDailyDriver(l, c.Daily)
			break
		case "single":
			instance, err = newSingleDriver(l, c.Singe)
			break
		case "stdout":
			instance = newStdoutDriver(l, c.StdOut)
			break
		default:
			err = errors.New("log channel is invalid")
		}

		if c.TeleHookConfig.Enable {
			addTelegramHook(l, c.TeleHookConfig)
		}
	})
	return instance, err
}
