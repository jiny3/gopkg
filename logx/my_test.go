package logx

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/jiny3/gopkg/filex"
)

func TestMyLogrus(t *testing.T) {
	file, _ := filex.FileOpen("logs/default.log")
	defer func() {
		file.Close()
		_ = filex.FileDelete("logs/default.log")
		_ = filex.DirDelete("logs")
	}()
	tests := []struct {
		name string
		want logrus.Level
	}{
		{
			name: "test logrus",
			want: logrus.DebugLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			MyLogrus()
			if !reflect.DeepEqual(logrus.GetLevel(), tt.want) {
				t.Errorf("MyLogrus() Level = %v, want %v", logrus.GetLevel().String(), tt.want.String())
			}
			logrus.Debug("test")
			wantLog := "test"
			// 读取日志文件内容
			logFile, err := os.Open("logs/default.log")
			if err != nil {
				t.Fatalf("failed opening file: %s", err)
			}
			defer logFile.Close()
			scanner := bufio.NewScanner(logFile)
			var line string
			if scanner.Scan() {
				line = scanner.Text()
			}
			if !strings.Contains(line, wantLog) {
				t.Errorf("MyLogrus() log content = %s, want %s", line, wantLog)
			}
			// 恢复标准输出
			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			line = buf.String()
			if !strings.Contains(line, wantLog) {
				t.Errorf("MyLogrus() stdout content = %s, want %s", line, wantLog+"\n")
			}
		})
	}
}
