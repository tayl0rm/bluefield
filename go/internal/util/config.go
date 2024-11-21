package util

import (
	v "github.com/spf13/viper"
)

type Config struct {
	ProjectID      string `mapstructure:"ProjectID"`
	BotToken       string `mapstructure:"BotToken"`
	Instance       string `mapstructure:"Instance"`
	Zone           string `mapstructure:"Zone"`
	ServerName     string `mapstructure:"ServerName"`
	ServerPassword string `mapstructure:"ServerPassword"`
}

func LoadConfig(path string) (config Config, err error) {
	v.AddConfigPath(path)
	v.SetConfigName("bot")
	v.SetConfigType("env")

	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err != nil {
		return
	}

	err = v.Unmarshal(&config)
	return
}
