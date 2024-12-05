package logx

import (
	"io"
	"os"

	"github.com/jiny3/gopkg/filex"
	"github.com/sirupsen/logrus"
)

var (
	Std  *logrus.Logger
	File *logrus.Logger
	All  *logrus.Logger

	// for old code
	Mystd  *logrus.Logger
	Myfile *logrus.Logger
	Myall  *logrus.Logger
)

func init() {
	stdInit()
	fileInit()
	allInit()

	Mystd = Std
	Myfile = File
	Myall = All
}

func stdInit() {
	Std = logrus.New()
	Std.SetOutput(os.Stdout)
	Std.SetLevel(logrus.DebugLevel)
	Std.SetFormatter(defaultFormatter)
}

func fileInit() {
	File = logrus.New()
	myLogWriter := []io.Writer{}
	var logConf struct {
		Writers []string
	}
	err := filex.ConfigRead("config", "log", &logConf)
	if err != nil {
		logConf = struct{ Writers []string }{
			Writers: []string{"default.log"},
		}
	}
	for _, path := range logConf.Writers {
		w, err := filex.FileOpen(path)
		if err != nil {
			logrus.Fatalf("open log file %s failed: %v", path, err)
		}
		myLogWriter = append(myLogWriter, w)
	}

	File.SetOutput(io.MultiWriter(myLogWriter...))
	File.SetLevel(logrus.TraceLevel)
	File.SetFormatter(defaultFormatter)
	File.SetReportCaller(true)
}

func allInit() {
	All = logrus.New()
	myLogWriter := []io.Writer{os.Stdout}
	var logConf struct {
		Writers []string
	}
	err := filex.ConfigRead("config", "log", &logConf)
	if err != nil {
		logConf = struct{ Writers []string }{
			Writers: []string{"default.log"},
		}
	}
	for _, path := range logConf.Writers {
		w, err := filex.FileOpen(path)
		if err != nil {
			logrus.Fatalf("open log file %s failed: %v", path, err)
		}
		myLogWriter = append(myLogWriter, w)
	}

	All.SetOutput(io.MultiWriter(myLogWriter...))
	All.SetLevel(logrus.TraceLevel)
	All.SetFormatter(defaultFormatter)
	All.SetReportCaller(true)
}
