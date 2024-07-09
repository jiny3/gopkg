package logx

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

func (f *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
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

func (f *myFormatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
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

func (f *myFormatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
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

func (f *myFormatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	fmt.Fprintf(b, " | %s:%v", field, entry.Data[field])
}
