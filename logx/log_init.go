package logx

import (
	"io"
	"os"

	"gitee.com/jiny1419/ucasnj-bbs/pkg/filex"
	"github.com/sirupsen/logrus"
)

func init() {
	myStdInit()
	myFileInit()
	myAllInit()
}

func myStdInit() {
	MyStd = logrus.New()
	MyStd.SetOutput(os.Stdout)
	MyStd.SetLevel(logrus.DebugLevel)
	MyStd.SetFormatter(&myFormatter{})
}

func myFileInit() {
	MyFile = logrus.New()
	myLogWriter := []io.Writer{}
	var logConf LogConfig
	err := filex.ReadConfig("config", "log", &logConf)
	if err != nil {
		logConf = LogConfig{
			Writers: []string{"default.log"},
		}
	}
	for _, path := range logConf.Writers {
		filex.CreateFile(path)
		w, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logrus.Fatalf("open log file %s failed: %v", path, err)
		}
		myLogWriter = append(myLogWriter, w)
	}

	MyFile.SetOutput(io.MultiWriter(myLogWriter...))
	MyFile.SetLevel(logrus.TraceLevel)
	MyFile.SetFormatter(privateMsgformatter)
	MyFile.SetReportCaller(true)
}

func myAllInit() {
	MyAll = logrus.New()
	myLogWriter := []io.Writer{os.Stdout}
	var logConf LogConfig
	err := filex.ReadConfig("config", "log", &logConf)
	if err != nil {
		logConf = LogConfig{
			Writers: []string{"default.log"},
		}
	}
	for _, path := range logConf.Writers {
		filex.CreateFile(path)
		w, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logrus.Fatalf("open log file %s failed: %v", path, err)
		}
		myLogWriter = append(myLogWriter, w)
	}

	MyAll.SetOutput(io.MultiWriter(myLogWriter...))
	MyAll.SetLevel(logrus.TraceLevel)
	MyAll.SetFormatter(privateMsgformatter)
	MyAll.SetReportCaller(true)
}
