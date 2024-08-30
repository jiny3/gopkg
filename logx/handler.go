package logx

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/JinY3/gopkg/filex"
	"github.com/sirupsen/logrus"
)

func DefaultConfig() struct {
	Level        logrus.Level
	Formatter    logrus.Formatter
	Output       io.Writer
	ReportCaller bool
} {
	myLogWriter := []io.Writer{os.Stdout}
	var logConf LogConfig
	err := filex.ReadConfig("config", "log", &logConf)
	if err != nil {
		logConf = LogConfig{
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
	return struct {
		Level        logrus.Level
		Formatter    logrus.Formatter
		Output       io.Writer
		ReportCaller bool
	}{
		Level: logrus.TraceLevel,
		Formatter: &myFormatter{
			Role:            "DEFAULT",
			TimestampFormat: "2006/01/02 - 15:04:05",
			CustomCallerFormatter: func(f *runtime.Frame) string {
				s := strings.Split(f.Function, ".")
				funcName := s[len(s)-1]
				return fmt.Sprintf("%s#%d:%s()", path.Base(f.File), f.Line, funcName)
			},
		},
		Output:       io.MultiWriter(myLogWriter...),
		ReportCaller: true,
	}
}
