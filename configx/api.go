package configx

import (
	"github.com/spf13/viper"
)

// set global viper config path,
// support ENV override: such as TEST_ENV_ONE override test.env.one
func Load(path string) error {
	return load(viper.GetViper(), path)
}

func Read[T any](path string, config *T) error {
	v := viper.New()
	err := load(v, path)
	if err != nil {
		return err
	}
	err = v.Unmarshal(config)
	if err != nil {
		return err
	}
	return nil
}
