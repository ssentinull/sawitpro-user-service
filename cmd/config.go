package main

import (
	"os"
	"time"

	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/joho/godotenv"
)

var (
	conf Config
)

type Config struct {
	Environment string
	Database    utils.DBOptions
	Auth        utils.AuthOptions
}

func loadConfig() (err error) {
	godotenv.Load()

	conf.Database.DSN = os.Getenv("DATABASE_DSN")
	conf.Auth.JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
	jwtExpiryDurationVar := os.Getenv("JWT_EXPIRY_DURATION")

	conf.Auth.JWTExpiryDuration, err = time.ParseDuration(jwtExpiryDurationVar)
	if err != nil {
		return err
	}

	return nil
}
