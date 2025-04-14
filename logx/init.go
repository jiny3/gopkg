package logx

import (
	"io"
	"os"

	"github.com/jiny3/gopkg/filex"
	"github.com/sirupsen/logrus"
)

// Init global logrus
func Init(level logrus.Level, path ...string) {
	logWriters := []io.Writer{os.Stdout}
	if len(path) == 0 {
		filex.FileDelete("logs/default.log")
		w, err := filex.FileOpen("logs/default.log")
		if err != nil {
			logrus.WithError(err).WithField("path", "logs/default.log").Fatal("open log file failed")
		}
		logWriters = append(logWriters, w)
	} else {
		for _, p := range path {
			filex.FileDelete(p)
			w, err := filex.FileOpen(p)
			if err != nil {
				logrus.WithError(err).WithField("path", p).Fatal("open log file failed")
			}
			logWriters = append(logWriters, w)
		}
	}
	logrus.SetOutput(io.MultiWriter(logWriters...))
	logrus.SetLevel(level)
	logrus.SetFormatter(defaultFormatter)
	logrus.SetReportCaller(true)
}
