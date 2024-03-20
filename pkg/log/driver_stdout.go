package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

type stdoutDriver struct {
	l *logrus.Logger
	c StdOutConfig
}

type StdOutConfig struct {
	Level string
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

func newStdoutDriver(l *logrus.Logger, c StdOutConfig) Contract {
	level := c.Level
	if level != "" {
		l.SetLevel(getLogLevel(level))
	}
	l.SetOutput(os.Stdout)
	return &stdoutDriver{
		l: l,
		c: c,
	}
}
