package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
)

type singleDriver struct {
	l *logrus.Logger
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
	fileName := viper.GetString("LOG_SINGLE_FILE_NAME")
	if fileName == "" {
		fileName = "app.log"
	}

	return fmt.Sprintf(
		"%s/%s",
		getLogsDir(),
		normalizedFilename(fileName),
	)
}

func newSingleDriver(l *logrus.Logger) (Contract, error) {
	level := viper.GetString("LOG_SINGLE_LOG_LEVEL")
	if level != "" {
		l.SetLevel(getLogLevel(level))
	}
	return &singleDriver{
		l: l,
	}, nil
}
