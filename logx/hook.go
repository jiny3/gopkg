package logx

import (
	"io"

	"github.com/sirupsen/logrus"
)

type hook struct {
	Writer    io.Writer
	Formatter logrus.Formatter
	Level     []logrus.Level
}

func (h *hook) Fire(entry *logrus.Entry) error {
	line, err := h.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = h.Writer.Write(line)
	if err != nil {
		return err
	}
	return nil
}

func (h *hook) Levels() []logrus.Level {
	return h.Level
}
