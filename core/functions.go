package core

import "github.com/spf13/viper"

func ReadViper() *viper.Viper {
	viperConfig := viper.New()
	viperConfig.SetConfigFile("config.yaml")
	_ = viperConfig.ReadInConfig()

	return viperConfig
}