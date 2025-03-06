package configx

import (
	"github.com/sirupsen/logrus"
)

func init() {
	err := Load("conf/config.yaml")
	if err != nil {
		logrus.WithField("path", "conf/config.yaml").Debug("default config not found")
	}
}
