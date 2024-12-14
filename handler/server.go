package handler

import (
	"github.com/SawitProRecruitment/UserService/usecase"
	"github.com/SawitProRecruitment/UserService/utils"
)

type Server struct {
	AuthUsecase usecase.AuthUsecaseInterface
	UserUsecase usecase.UserUsecaseInterface
	AuthUtil    utils.AuthInterface
}

type NewServerOptions struct {
	AuthUsecase usecase.AuthUsecaseInterface
	UserUsecase usecase.UserUsecaseInterface
	AuthUtil    utils.AuthInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		AuthUsecase: opts.AuthUsecase,
		UserUsecase: opts.UserUsecase,
		AuthUtil:    opts.AuthUtil,
	}
}
