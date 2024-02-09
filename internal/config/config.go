package config

import (
	"github.com/spf13/viper"
)

type DBConfig struct {
	FileName string `mapstructure:"file"`
}

type Config struct {
	db DBConfig `mapstructure:"db"`
}

// GetConfig returns data from config file
func GetConfig(confPath string) (Config, error) {
	var conf Config

	viper.SetConfigName("transactions-summary")
	viper.AddConfigPath(confPath)

	err := viper.ReadInConfig()
	if err != nil {
		return conf, err
	}

	err = viper.Unmarshal(&conf)

	return conf, err
}
