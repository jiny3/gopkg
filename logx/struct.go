package logx

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var MyStd *logrus.Logger
var MyFile *logrus.Logger
var MyAll *logrus.Logger

type LogConfig struct {
	Writers []string
}

type myFormatter struct {
	Role                  string
	TimestampFormat       string
	CustomCallerFormatter func(*runtime.Frame) string
}

var privateMsgformatter = &myFormatter{
	Role:            "DEFAULT",
	TimestampFormat: "2006/01/02 - 15:04:05",
	CustomCallerFormatter: func(f *runtime.Frame) string {
		s := strings.Split(f.Function, ".")
		funcName := s[len(s)-1]
		return fmt.Sprintf("%s#%d:%s()", path.Base(f.File), f.Line, funcName)
	},
}
