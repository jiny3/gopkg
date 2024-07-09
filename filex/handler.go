package filex

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

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

func GetJson[T any](file string, entrys *[]T) error {
	openFile, err := os.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer openFile.Close()

	decoder := json.NewDecoder(openFile)
	for {
		var entry T
		if err := decoder.Decode(&entry); err != nil {
			break
		}
		*entrys = append(*entrys, entry)
	}

	return nil
}

func InsertJson[T any](file string, entry T) error {
	// T to json
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	// open file
	openFile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer openFile.Close()
	// write json
	if _, err := openFile.Write(jsonData); err != nil {
		return err
	}
	return nil
}

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

func CreateFile(path string) error {
	dir := filepath.Dir(path)
	if err := CreateDir(dir); err != nil {
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
