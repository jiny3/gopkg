package hookx

import (
	"github.com/jiny3/gopkg/configx"
	"github.com/jiny3/gopkg/logx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var DefaultHook = func() {
	err := configx.Load("config/config.yaml")
	if err != nil {
		logx.Init(logrus.DebugLevel)
		logrus.WithError(err).Error("config load failed")
		return
	}
	logPath, level := viper.GetString("log.path"), viper.GetString("log.level")
	_level, err := logrus.ParseLevel(level)
	if err != nil {
		_level = logrus.InfoLevel
	}
	if logPath == "" {
		logx.Init(_level)
	} else {
		logx.Init(_level, logPath)
	}
}
