package logx

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/jiny3/gopkg/configx"
	"github.com/jiny3/gopkg/filex"
)

type config struct {
	Level        logrus.Level
	Formatter    logrus.Formatter
	Output       io.Writer
	ReportCaller bool
}

// Deprecated: This function will be removed in a future version.
// config is unnecessary, use NewLogger or InitLogrus instead.
func DefaultConfig() config {
	logWriter := []io.Writer{os.Stdout}
	var logConf struct {
		Writers []string
	}
	err := configx.Read("config/log.yaml", &logConf)
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
		Formatter:    defaultFormatter(),
		Output:       io.MultiWriter(logWriter...),
		ReportCaller: true,
	}
}

func initLogger(logger *logrus.Logger, opts ...Option) {
	if logger == nil {
		logger = logrus.New()
	}
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(defaultFormatter())
	logger.SetReportCaller(true)
	for _, opt := range opts {
		opt(logger)
	}
}

func NewLogger(opts ...Option) *logrus.Logger {
	logger := logrus.New()
	initLogger(logger, opts...)
	return logger
}

// InitLogrus global logrus
func InitLogrus(opts ...Option) {
	initLogger(logrus.StandardLogger(), opts...)
}
