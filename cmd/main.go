package main

import (
	"log"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/labstack/echo/v4"
)

func init() {
	if err := loadConfig(); err != nil {
		log.Fatal("error loading .env: ", err)
	}
}

func main() {
	e := echo.New()
	server := newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":" + conf.ServicePort))
}

func newServer() *handler.Server {
	repo := repository.NewRepository(repository.NewRepositoryOptions{Dsn: conf.DatabaseURL})
	opts := handler.NewServerOptions{Repository: repo}

	return handler.NewServer(opts)
}
