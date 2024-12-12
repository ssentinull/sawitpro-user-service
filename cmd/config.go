package main

import (
	"github.com/spf13/viper"
)

var (
	conf Config
)

type Config struct {
	ServicePort string
	DatabaseURL string
}

func loadConfig() (err error) {
	defer func() {
		if errRecov := recover(); errRecov != nil {
			err = errRecov.(error)
		}
	}()

	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	conf.ServicePort = viper.GetString("SERVICE_PORT")
	conf.DatabaseURL = viper.GetString("DATABASE_URL")

	return nil
}
