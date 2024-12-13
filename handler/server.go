package handler

import (
	"github.com/SawitProRecruitment/UserService/usecase"
)

type Server struct {
	AuthUsecase usecase.AuthUsecaseInterface
	UserUsecase usecase.UserUsecaseInterface
}

type NewServerOptions struct {
	AuthUsecase usecase.AuthUsecaseInterface
	UserUsecase usecase.UserUsecaseInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		AuthUsecase: opts.AuthUsecase,
		UserUsecase: opts.UserUsecase,
	}
}
