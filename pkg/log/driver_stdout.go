package log

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type stdoutDriver struct {
	l *logrus.Logger
}

func (d stdoutDriver) Debug(args ...interface{}) {
	d.getLog().Debug(args)
}

func (d stdoutDriver) Info(args ...interface{}) {
	d.getLog().Info(args)
}

func (d stdoutDriver) Warn(args ...interface{}) {
	d.getLog().Warn(args)
}

func (d stdoutDriver) Error(args ...interface{}) {
	d.getLog().Error(args)
}

func (d stdoutDriver) Fatal(args ...interface{}) {
	d.getLog().Fatal("fatal", args)
}

func (d stdoutDriver) Fatalf(format string, args ...interface{}) {
	d.getLog().Fatalf(format, args)
}
func (d stdoutDriver) getLog() *logrus.Entry {
	return d.l.WithField("file", getCalledFile(3))
}

func newStdoutDriver(l *logrus.Logger) Contract {
	level := viper.GetString("LOG_STDOUT_LOG_LEVEL")
	if level != "" {
		l.SetLevel(getLogLevel(level))
	}
	l.SetOutput(os.Stdout)
	return &stdoutDriver{
		l: l,
	}
}
