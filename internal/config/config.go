package config

import (
	"github.com/spf13/viper"

	"myservice.com/transactions/internal/clients/sender"
)

type DBConfig struct {
	FileName string `mapstructure:"file"`
}

type Config struct {
	db          DBConfig      `mapstructure:"db"`
	EmailSender sender.Config `mapstructure:"sendgrid"`
}

// GetConfig returns data from config file
func GetConfig(confPath string) (Config, error) {
	var conf Config

	viper.SetConfigName("transactions-summary-service")
	viper.AddConfigPath(confPath)

	err := viper.ReadInConfig()
	if err != nil {
		return conf, err
	}

	err = viper.Unmarshal(&conf)

	return conf, err
}
