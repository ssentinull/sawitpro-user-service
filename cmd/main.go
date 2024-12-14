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
	server, err := newServer()
	if err != nil {
		panic(err)
	}

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":" + conf.ServicePort))
}

func newServer() (*handler.Server, error) {
	DB, err := utils.InitDB(conf.Database)
	if err != nil {
		return nil, err
	}

	auth, err := utils.InitAuth(conf.Auth)
	if err != nil {
		return nil, err
	}

	userRepo := repository.NewUserRepository(repository.UserRepositoryOptions{DB: DB})
	authUsecase := usecase.NewAuthUsecase(usecase.AuthUsecaseOptions{
		UserRepository: userRepo,
		AuthUtil:       auth,
	})

	userUsecase := usecase.NewUserUsecase(usecase.UserUsecaseOptions{UserRepository: userRepo})
	opts := handler.NewServerOptions{
		AuthUsecase: authUsecase,
		UserUsecase: userUsecase,
		AuthUtil:    auth,
	}

	return handler.NewServer(opts), nil
}
