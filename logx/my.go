package logx

import (
	"github.com/sirupsen/logrus"
)

func MyLogrus() {
	InitLogrus(WithLevel(logrus.DebugLevel), WithAllText("logs/default.log"))
}
