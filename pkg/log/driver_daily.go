package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

type dailyDriver struct {
	l *logrus.Logger
	c DailyConfig
}

type DailyConfig struct {
	FileName string
	Level    string
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
	return fmt.Sprintf(
		"%s/%s",
		getLogsDir(),
		d.getLogFileName(),
	)
}

func (d dailyDriver) getLogFileName() string {
	fileName := d.c.FileName
	if fileName == "" {
		fileName = "app.log"
	}
	return fmt.Sprintf(
		"%s-%s.log",
		normalizedFilename(fileName),
		time.Now().Format("2006-01-02"),
	)
}

func newDailyDriver(l *logrus.Logger, c DailyConfig) (Contract, error) {
	level := c.Level
	if level != "" {
		l.SetLevel(getLogLevel(level))
	}
	return &dailyDriver{
		l: l,
		c: c,
	}, nil
}
