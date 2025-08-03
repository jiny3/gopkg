package hookx

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/jiny3/gopkg/configx"
	"github.com/jiny3/gopkg/logx"
)

var WithDefault = func() {
	err := configx.Load("config/config.yaml")
	if err != nil {
		logx.MyLogrus()
		logrus.WithError(err).Error("config load failed")
		return
	}
	logPath, level := viper.GetString("log.path"), viper.GetString("log.level")
	_level, err := logrus.ParseLevel(level)
	if err != nil {
		_level = logrus.InfoLevel
	}
	if logPath == "" {
		logx.InitLogrus(logx.WithLevel(_level), logx.WithAllText("logs/default.log"))
	} else {
		logx.InitLogrus(logx.WithLevel(_level), logx.WithAllText(logPath))
	}
}

// enable pprof for debug, default addr :6060
var WithPprof = func() {
	go http.ListenAndServe(":6060", nil)
}

// enable prometheus for monitor, default addr :8123
var WithPrometheus = func() {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8123", nil)
}
