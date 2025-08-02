package configx

import (
	"strings"

	"github.com/jiny3/gopkg/filex"
	"github.com/spf13/viper"
)

// set input viper config
func load(v *viper.Viper, path string) error {
	args := filex.Parse(path)
	v.SetConfigName(args.Name)
	v.SetConfigType(args.Type)
	v.AddConfigPath(args.Dir)

	// env 覆盖
	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	return v.ReadInConfig()
}
