package logx

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/jiny3/gopkg/filex"
	"github.com/sirupsen/logrus"
)

type config struct {
	Level        logrus.Level
	Formatter    logrus.Formatter
	Output       io.Writer
	ReportCaller bool
}

func DefaultConfig() config {
	logWriter := []io.Writer{os.Stdout}
	var logConf struct {
		Writers []string
	}
	err := filex.ReadConfig("config", "log", &logConf)
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
		logWriter = append(logWriter, w)
	}
	return config{
		Level:        logrus.TraceLevel,
		Formatter:    DebugFormatter(),
		Output:       io.MultiWriter(logWriter...),
		ReportCaller: true,
	}
}

func DebugFormatter() logrus.Formatter {
	return &Formatter{
		Role:            "DEFAULT",
		TimestampFormat: "2006/01/02 - 15:04:05",
		CustomCallerFormatter: func(f *runtime.Frame) string {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return fmt.Sprintf("%s#%d:%s()", path.Base(f.File), f.Line, funcName)
		},
	}
}
