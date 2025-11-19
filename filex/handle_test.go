package filex

import (
	"os"
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

func TestFileCreate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path    string
		wantErr bool
	}{
		{
			name:    "create file in new directory",
			path:    "path/to/new/file.txt",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := FileCreate(tt.path)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("FileCreate() failed: %v", gotErr)
				}
				return
			}
			defer os.RemoveAll("path") // clean up after test
			if tt.wantErr {
				t.Fatal("FileCreate() succeeded unexpectedly")
			}
		})
	}
}
