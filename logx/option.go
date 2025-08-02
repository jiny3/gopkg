package logx

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/jiny3/gopkg/filex"
)

type Option func(*logrus.Logger)

func WithLevel(level logrus.Level) Option {
	return func(l *logrus.Logger) {
		l.SetLevel(level)
	}
}

func WithFormatter(formatter logrus.Formatter) Option {
	return func(l *logrus.Logger) {
		l.SetFormatter(formatter)
	}
}

func WithCallerReport(report bool) Option {
	return func(l *logrus.Logger) {
		l.SetReportCaller(report)
	}
}

// WithOutput sets the output for the logger.
// std true will write to os.Stdout.
// logPath will write to the specified log files.
func WithOutput(std bool, logPath ...string) Option {
	return func(l *logrus.Logger) {
		logWriters := []io.Writer{}
		if std {
			logWriters = append(logWriters, os.Stdout)
		}
		for _, p := range logPath {
			w, err := filex.FileOpen(p)
			if err != nil {
				logrus.WithError(err).WithField("path", p).Fatal("open log file failed")
			}
			logWriters = append(logWriters, w)
		}
		l.SetOutput(io.MultiWriter(logWriters...))
	}
}
