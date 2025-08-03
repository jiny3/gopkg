package logx

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// Deprecated: This type will be removed in a future version.
// Use TextFormatter instead.
type Formatter = TextFormatter

type TextFormatter struct {
	Role                  string
	TimestampFormat       string
	CustomCallerFormatter func(*runtime.Frame) string
	ForceColors           bool
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 创建一个缓冲区
	b := &bytes.Buffer{}

	role := f.Role
	timestamp := entry.Time.Format(f.TimestampFormat)
	level := strings.ToUpper(entry.Level.String())

	// 自定义日志头
	if f.ForceColors {
		switch entry.Level {
		case logrus.TraceLevel, logrus.DebugLevel:
			b.WriteString(strings.Join([]string{
				"[" + role + "] ",
				colorize(timestamp, int(color.FgGreen)),
				"  ",
				colorize(level, int(color.BgBlue)),
			}, ""))
		case logrus.WarnLevel:
			b.WriteString(strings.Join([]string{
				"[" + role + "] ",
				colorize(timestamp, int(color.FgGreen)),
				"  ",
				colorize(level, int(color.BgYellow)),
			}, ""))
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			b.WriteString(strings.Join([]string{
				"[" + role + "] ",
				colorize(timestamp, int(color.FgGreen)),
				"  ",
				colorize(level, int(color.BgRed)),
			}, ""))
		default:
			b.WriteString(strings.Join([]string{
				"[" + role + "] ",
				colorize(timestamp, int(color.FgGreen)),
				"  ",
				colorize(level, int(color.BgGreen)),
			}, ""))
		}
	} else {
		b.WriteString(strings.Join([]string{
			"[" + role + "] ",
			timestamp,
			" | ",
			level,
		}, ""))
	}

	// 写入caller
	f.writeCaller(b, entry)

	// 写入fields
	f.writeFields(b, entry)

	// 写入日志内容
	if f.ForceColors {
		switch entry.Level {
		case logrus.TraceLevel, logrus.DebugLevel:
			b.WriteString(strings.Join([]string{
				"    ",
				colorize(entry.Message, int(color.FgHiBlue)),
			}, ""))
		case logrus.WarnLevel:
			b.WriteString(strings.Join([]string{
				"    ",
				colorize(entry.Message, int(color.FgHiYellow)),
			}, ""))
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			b.WriteString(strings.Join([]string{
				"    ",
				colorize(entry.Message, int(color.FgHiRed)),
			}, ""))
		default:
			b.WriteString(strings.Join([]string{
				"    ",
				colorize(entry.Message, int(color.FgHiGreen)),
			}, ""))
		}
	} else {
		b.WriteString(strings.Join([]string{
			"    ",
			entry.Message,
		}, ""))
	}

	// 写入换行
	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *TextFormatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		if f.ForceColors {
			b.WriteString("  ")
		} else {
			b.WriteString(" | ")
		}
		if f.CustomCallerFormatter != nil {
			fmt.Fprint(b, f.CustomCallerFormatter(entry.Caller))
		} else {
			b.WriteString(strings.Join([]string{
				entry.Caller.File,
				"#",
				strconv.Itoa(entry.Caller.Line),
				":",
				entry.Caller.Function,
			}, ""))
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
	if f.ForceColors {
		fmt.Fprintf(b, "  %s:%s", colorize(field, int(color.FgHiMagenta)), colorize(fmt.Sprintf("%v", entry.Data[field]), int(color.FgMagenta)))
	} else {
		fmt.Fprintf(b, " | %s:%v", field, entry.Data[field])
	}
}

func colorize(s string, color int) string {
	return strings.Join([]string{"\x1b[", strconv.Itoa(color), "m", s, "\x1b[0m"}, "")
}

func defaultFormatter(isColor bool) *TextFormatter {
	if isColor {
		return &TextFormatter{
			Role:            "DEFAULT",
			TimestampFormat: "2006/01/02 - 15:04:05",
			CustomCallerFormatter: func(f *runtime.Frame) string {
				s := strings.Split(f.Function, ".")
				funcName := s[len(s)-1]
				return colorize(strings.Join([]string{path.Base(f.File), "#", strconv.Itoa(f.Line), ":", funcName, "()"}, ""), int(color.Underline))
			},
			ForceColors: true,
		}
	}

	return &TextFormatter{
		Role:            "DEFAULT",
		TimestampFormat: "2006/01/02 - 15:04:05",
		CustomCallerFormatter: func(f *runtime.Frame) string {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return strings.Join([]string{path.Base(f.File), "#", strconv.Itoa(f.Line), ":", funcName, "()"}, "")
		},
		ForceColors: true,
	}
}
