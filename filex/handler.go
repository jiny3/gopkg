package filex

import (
	"os"
	"path/filepath"
	"strings"
)

// parse path to dir, name, type
func Parse(path string) struct {
	Dir, Name, Type string
} {
	dir, file := path[:strings.LastIndex(path, "/")], path[strings.LastIndex(path, "/")+1:]
	name, typ := file[:strings.LastIndex(file, ".")], file[strings.LastIndex(file, ".")+1:]
	return struct {
		Dir, Name, Type string
	}{dir, name, typ}
}

func DirCreate(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// if not exist, create it
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func DirDelete(path string) error {
	// if not exist, pass
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return nil
}

func FileCreate(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if err := DirCreate(dir); err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
}

func FileDelete(path string) error {
	// if not exist, pass
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}

func FileOpen(path string) (*os.File, error) {
	return FileCreate(path)
}
