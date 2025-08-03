package logx

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

// Deprecated: This type will be removed in a future version.
// Use TextFormatter instead.
type Formatter = TextFormatter

type TextFormatter struct {
	logrus.TextFormatter
	Role                  string
	TimestampFormat       string
	CustomCallerFormatter func(*runtime.Frame) string
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 创建一个缓冲区
	b := &bytes.Buffer{}

	role := f.Role
	timestamp := entry.Time.Format(f.TimestampFormat)
	level := strings.ToUpper(entry.Level.String())

	// 自定义日志头
	b.WriteString(fmt.Sprintf("[%s] %s | %s", role, timestamp, level))

	// 写入caller
	f.writeCaller(b, entry)

	// 写入fields
	f.writeFields(b, entry)

	// 写入日志内容
	b.WriteString(fmt.Sprintf("    \"%s\"", entry.Message))

	// 写入换行
	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *TextFormatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		b.WriteString(" | ")
		if f.CustomCallerFormatter != nil {
			fmt.Fprint(b, f.CustomCallerFormatter(entry.Caller))
		} else {
			fmt.Fprintf(
				b,
				"%s#%d:%s()",
				entry.Caller.File,
				entry.Caller.Line,
				entry.Caller.Function,
			)
		}
	}
}

func (f *TextFormatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *TextFormatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	fmt.Fprintf(b, " | %s:%v", field, entry.Data[field])
}

func defaultFormatter() *TextFormatter {
	return &TextFormatter{
		Role:            "DEFAULT",
		TimestampFormat: "2006/01/02 - 15:04:05",
		CustomCallerFormatter: func(f *runtime.Frame) string {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return fmt.Sprintf("%s#%d:%s()", path.Base(f.File), f.Line, funcName)
		},
		TextFormatter: logrus.TextFormatter{
			ForceColors: true,
		},
	}
}
