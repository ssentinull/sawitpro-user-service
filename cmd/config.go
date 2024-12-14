package main

import (
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/spf13/viper"
)

var (
	conf Config
)

type Config struct {
	ServicePort string
	Database    utils.DBOptions
	Auth        utils.AuthOptions
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
	conf.Database.DSN = viper.GetString("DATABASE_DSN")
	conf.Auth.JWTExpiryDuration = viper.GetDuration("JWT_EXPIRY_DURATION")
	conf.Auth.JWTSecretKey = viper.GetString("JWT_SECRET_KEY")

	return nil
}
