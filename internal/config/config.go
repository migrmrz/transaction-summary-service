package config

import (
	"gopkg.in/yaml.v3"

	"myservice.com/transactions/internal/clients/sender"
)

type DBConfig struct {
	FileName string `mapstructure:"file"`
}

type Config struct {
	TransactionsFile string        `mapstructure:"transactions-file" yaml:"transactions-file"`
	EmailSender      sender.Config `mapstructure:"sendgrid" yaml:"sendgrid"`
}

// GetConfig returns data from config file
func GetConfigFromS3(data []byte) (Config, error) {
	var conf Config

	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
