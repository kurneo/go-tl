package logger

import (
	"fmt"
	"github.com/kurneo/go-template/config"
	"github.com/sirupsen/logrus"
	"log"
	"path/filepath"
	"strings"
	"time"
)

type dailyDriver struct {
	l *logrus.Logger
	c config.LogChannel
}

func (d dailyDriver) Debug(args ...interface{}) {
	d.getLog().Debug(args)
}

func (d dailyDriver) Info(args ...interface{}) {
	d.getLog().Info(args)
}

func (d dailyDriver) Warn(args ...interface{}) {
	d.getLog().Warn(args)
}

func (d dailyDriver) Error(args ...interface{}) {
	d.getLog().Error(args)
}

func (d dailyDriver) Fatal(args ...interface{}) {
	d.getLog().Fatal("fatal", args)
}

func (d dailyDriver) Fatalf(format string, args ...interface{}) {
	d.getLog().Fatalf(format, args)
}

func (d dailyDriver) getLog() *logrus.Entry {
	d.prepareLogFile()
	return d.l.WithField("file", getCalledFile(3))
}

func (d dailyDriver) prepareLogFile() {
	if file == nil {
		f, err := createLogFile(d.getLogFilePath())
		if err != nil {
			log.Fatalf("Log error: cannot create new log file %s", err)
			return
		}
		file = f
	} else {
		if file.Name() == d.getLogFileName() {
			return
		}

		f, err := createLogFile(d.getLogFilePath())
		if err != nil {
			log.Fatalf("Log error: cannot create new log file %s", err)
			return
		}

		if file != nil {
			err = file.Close()
			if err != nil {
				log.Fatalf("Log error: cannot close old log file %s", err)
			}
		}
		file = f
	}
	d.l.SetOutput(file)
}

func (d dailyDriver) getLogFilePath() string {
	root := d.c.Get("root")
	if root == "" {
		root = "storage/logs"
	}
	return fmt.Sprintf(
		"%s/%s",
		root,
		d.getLogFileName(),
	)
}

func (d dailyDriver) getLogFileName() string {
	fileName := d.c.Get("file_name")
	if fileName == "" {
		fileName = "app.log"
	}
	return fmt.Sprintf(
		"%s-%s.log",
		strings.TrimSuffix(fileName, filepath.Ext(fileName)),
		time.Now().Format("2006-01-02"),
	)
}

func newDailyDriver(c config.LogChannel, l *logrus.Logger) (Contract, error) {
	if c.Get("level") != "" {
		l.SetLevel(getLogLevel(c.Get("level")))
	}
	return &dailyDriver{
		l: l,
		c: c,
	}, nil
}
