package logx

import (
	"github.com/sirupsen/logrus"
)

func MyLogrus() {
	InitLogrus(WithLevel(logrus.DebugLevel), WithOutput(true, "logs/default.log"))
}
