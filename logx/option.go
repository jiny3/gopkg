package logx

import (
	"io"

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
		if formatter == nil {
			return
		}
		l.SetFormatter(formatter)
	}
}

func WithCallerReport(report bool) Option {
	return func(l *logrus.Logger) {
		l.SetReportCaller(report)
	}
}

// WithAllText sets the output for the logger.
// level same as logrus.StandardLogger()
// logPath will write to the specified log files.
func WithAllText(logPath ...string) Option {
	return func(l *logrus.Logger) {
		if len(logPath) == 0 {
			return
		}
		curLevel := logrus.GetLevel()
		levels := []logrus.Level{}
		for _, level := range logrus.AllLevels {
			if level <= curLevel {
				levels = append(levels, level)
			}
		}
		logWriters := []io.Writer{}
		for _, p := range logPath {
			w, err := filex.FileOpen(p)
			if err != nil {
				logrus.WithError(err).WithField("path", p).Fatal("open log file failed")
			}
			logWriters = append(logWriters, w)
		}
		f := defaultFormatter()
		f.ForceColors = false
		hook := &hook{
			Writer:    io.MultiWriter(logWriters...),
			Formatter: f,
			Level:     levels,
		}
		l.AddHook(hook)
	}
}

// WithOpsJSON sets the logger to output JSON format for operational logs.
// It writes to the specified log files and includes only warning, error, fatal, and panic levels.
func WithOpsJSON(logPath ...string) Option {
	return func(l *logrus.Logger) {
		if len(logPath) == 0 {
			return
		}
		levels := []logrus.Level{
			logrus.WarnLevel,
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		}
		logWriters := []io.Writer{}
		for _, p := range logPath {
			w, err := filex.FileOpen(p)
			if err != nil {
				logrus.WithError(err).WithField("path", p).Fatal("open log file failed")
			}
			logWriters = append(logWriters, w)
		}
		hook := &hook{
			Writer:    io.MultiWriter(logWriters...),
			Formatter: &logrus.JSONFormatter{},
			Level:     levels,
		}
		l.AddHook(hook)
	}
}

// WithDevJSON sets the logger to output JSON format for development logs.
// It writes to the specified log files and includes trace and debug levels.
func WithDevJSON(logPath ...string) Option {
	return func(l *logrus.Logger) {
		if len(logPath) == 0 {
			return
		}
		levels := []logrus.Level{
			logrus.TraceLevel,
			logrus.DebugLevel,
		}
		logWriters := []io.Writer{}
		for _, p := range logPath {
			w, err := filex.FileOpen(p)
			if err != nil {
				logrus.WithError(err).WithField("path", p).Fatal("open log file failed")
			}
			logWriters = append(logWriters, w)
		}
		hook := &hook{
			Writer:    io.MultiWriter(logWriters...),
			Formatter: &logrus.JSONFormatter{},
			Level:     levels,
		}
		l.AddHook(hook)
	}
}
