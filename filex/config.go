package filex

import "github.com/spf13/viper"

// Deprecated: This function will be removed in a future version.
// Use configx.Read instead
func ConfigRead[T any](configParentPath string, configName string, config *T) error {
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
