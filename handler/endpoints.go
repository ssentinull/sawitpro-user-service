package handler

import (
	"fmt"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/labstack/echo/v4"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {

	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) AuthLogin(ctx echo.Context) error {
	req := generated.AuthLoginJSONRequestBody{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Success: false,
			Message: "Invalid Input.",
		})
	}

	user, jwt, err := s.AuthUsecase.LoginUser(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	resp := generated.AuthLoginResponse{
		Success: true,
		Message: "successfully logged-in user",
		Data: &struct {
			Id  int    "json:\"id\""
			Jwt string "json:\"jwt\""
		}{
			Id:  int(user.Id),
			Jwt: jwt,
		},
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) RegisterUser(ctx echo.Context) error {
	req := generated.RegisterUserJSONRequestBody{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Success: false,
			Message: "Invalid Input.",
		})
	}

	if isPayloadValid, errorMessage := utils.IsRegisterUserPayloadValid(req); !isPayloadValid {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Success: false,
			Message: errorMessage,
		})
	}

	// TODO: implement stacktrace
	user, err := s.UserUsecase.CreateUser(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	resp := generated.RegisterUserResponse{
		Success: true,
		Message: "successfully created user",
		Data: &struct {
			Id int "json:\"id\""
		}{
			Id: int(user.Id),
		},
	}

	return ctx.JSON(http.StatusCreated, resp)
}
