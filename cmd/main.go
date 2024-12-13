package main

import (
	"log"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/usecase"
	"github.com/SawitProRecruitment/UserService/utils"

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
	auth := utils.InitAuth(conf.Auth)

	repo := repository.NewRepository(repository.NewRepositoryOptions{Dsn: conf.DatabaseURL})
	authUsecase := usecase.NewAuthUsecase(usecase.AuthUsecaseOptions{
		Repository: repo,
		Auth:       auth,
	})

	userUsecase := usecase.NewUserUsecase(usecase.UserUsecaseOptions{Repository: repo})
	opts := handler.NewServerOptions{
		AuthUsecase: authUsecase,
		UserUsecase: userUsecase,
	}

	return handler.NewServer(opts)
}
