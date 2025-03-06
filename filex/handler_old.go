package filex

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Deprecated: This function will be removed in a future version.
// Use DirCreate instead
func CreateDir(path string) error {
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

// Deprecated: This function will be removed in a future version.
// Use FileCreate instead
func CreateFile(path string) error {
	dir := filepath.Dir(path)
	if err := DirCreate(dir); err != nil {
		return err
	}
	// if not exist, create it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	} else if err != nil {
		return err
	}
	return nil
}

// Deprecated: This function will be removed in a future version.
// Use FileDelete instead
func DeleteFile(path string) error {
	// if not exist, pass
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}

// Deprecated: This function will be removed in a future version.
// Use configx.Read instead
func ReadConfig[T any](configParentPath string, configName string, config *T) error {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configParentPath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(config)
	if err != nil {
		return err
	}
	return nil
}
