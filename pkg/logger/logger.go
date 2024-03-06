package logger

import (
	"errors"
	"github.com/kurneo/go-template/config"
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

func New(cfg config.Log) (Contract, error) {
	var err error
	var c map[string]string
	if _, ok := cfg.Channels[cfg.Default]; ok {
		c = cfg.Channels[cfg.Default]
	}
	once.Do(func() {
		level := getLogLevel(cfg.Level)
		l := newLogrus(level)
		setupHook(cfg, l)
		switch cfg.Default {
		case "daily":
			instance, err = newDailyDriver(c, l)
			break
		case "single":
			instance, err = newSingleDriver(c, l)
			break
		case "stdout":
			instance = newStdoutDriver(c, l)
			break
		default:
			err = errors.New("invalid log channel: " + cfg.Default)
		}
	})
	return instance, err
}
