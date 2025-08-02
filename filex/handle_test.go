package filex

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want struct{ Dir, Name, Type string }
	}{
		{
			name: "test normal",
			args: args{path: "path/to/file.txt"},
			want: struct{ Dir, Name, Type string }{
				Dir:  "path/to",
				Name: "file",
				Type: "txt",
			},
		},
		{
			name: "test with no extension",
			args: args{path: "path/to/file"},
			want: struct{ Dir, Name, Type string }{
				Dir:  "path/to",
				Name: "file",
				Type: "",
			},
		},
		{
			name: "test with multiple dots",
			args: args{path: "path/to/file.name.txt"},
			want: struct{ Dir, Name, Type string }{
				Dir:  "path/to",
				Name: "file.name",
				Type: "txt",
			},
		},
		{
			name: "test without directory",
			args: args{path: "file.txt"},
			want: struct{ Dir, Name, Type string }{
				Dir:  ".",
				Name: "file",
				Type: "txt",
			},
		},
		{
			name: "test with empty path",
			args: args{path: ""},
			want: struct{ Dir, Name, Type string }{
				Dir:  ".",
				Name: "",
				Type: "",
			},
		},
		{
			name: "test with only directory",
			args: args{path: "path/to/"},
			want: struct{ Dir, Name, Type string }{
				Dir:  "path/to",
				Name: "",
				Type: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
