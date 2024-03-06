package logger

import (
	"fmt"
	"github.com/kurneo/go-template/config"
	"github.com/sirupsen/logrus"
	"log"
)

type singleDriver struct {
	l *logrus.Logger
	c config.LogChannel
}

func (d singleDriver) Debug(args ...interface{}) {
	d.getLog().Debug(args)
}

func (d singleDriver) Info(args ...interface{}) {
	d.getLog().Info(args)
}

func (d singleDriver) Warn(args ...interface{}) {
	d.getLog().Warn(args)
}

func (d singleDriver) Error(args ...interface{}) {
	d.getLog().Error(args)
}

func (d singleDriver) Fatal(args ...interface{}) {
	d.getLog().Fatal("fatal", args)
}

func (d singleDriver) Fatalf(format string, args ...interface{}) {
	d.getLog().Fatalf(format, args)
}
func (d singleDriver) getLog() *logrus.Entry {
	d.prepareLogFile()
	return d.l.WithField("file", getCalledFile(3))
}

func (d singleDriver) prepareLogFile() {
	if file == nil {
		f, err := createLogFile(d.getLogFilePath())
		if err != nil {
			log.Fatalf("Log error: cannot create new log file %s", err)
			return
		}
		file = f
		d.l.SetOutput(file)
	}
}

func (d singleDriver) getLogFilePath() string {
	root := d.c.Get("root")
	if root == "" {
		root = "storage/logs"
	}

	fileName := d.c.Get("file_name")
	if fileName == "" {
		fileName = "app.log"
	}

	return fmt.Sprintf(
		"%s/%s",
		root,
		fileName,
	)
}

func newSingleDriver(c config.LogChannel, l *logrus.Logger) (Contract, error) {
	if c.Get("level") != "" {
		l.SetLevel(getLogLevel(c.Get("level")))
	}
	return &singleDriver{
		l: l,
		c: c,
	}, nil
}
