package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/utils"

	"github.com/labstack/echo/v4"
)

func (s *Server) AuthLogin(ctx echo.Context) error {
	req := generated.AuthLoginJSONRequestBody{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Success: false,
			Message: "Invalid Input.",
		})
	}

	if isPayloadValid, errorMessage := utils.IsAuthLoginPayloadValid(req); !isPayloadValid {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Success: false,
			Message: errorMessage,
		})
	}

	user, jwt, err := s.AuthUsecase.LoginUser(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(int(utils.GetCode(err)), generated.ErrorResponse{
			Success: false,
			Message: utils.GetMessage(err),
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

	user, err := s.UserUsecase.CreateUser(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(int(utils.GetCode(err)), generated.ErrorResponse{
			Success: false,
			Message: utils.GetMessage(err),
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

func (s *Server) GetUserProfile(ctx echo.Context) error {
	tokenStr := ctx.Request().Header.Get("authorization")
	idx := strings.Index(tokenStr, " ")
	if tokenStr == "" || idx < 0 {
		return errors.New("invalid jwt token")
	}

	tokenStr = tokenStr[idx+1:]
	if err := s.AuthUtil.ValidateJWTToken(tokenStr); err != nil {
		return err
	}

	userId, err := s.AuthUtil.GetUserId(tokenStr)
	if err != nil {
		return err
	}

	user, err := s.UserUsecase.GetUserProfile(ctx.Request().Context(), userId)
	if err != nil {
		return err
	}

	resp := generated.GetUserProfileResponse{
		Success: true,
		Message: "successfully get user profile",
		Data: &struct {
			Id          int    "json:\"id\""
			FullName    string "json:\"full_name\""
			PhoneNumber string "json:\"phone_number\""
		}{
			Id:          int(user.Id),
			FullName:    user.FullName,
			PhoneNumber: user.PhoneNumber,
		},
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateUserProfile(ctx echo.Context) error {
	tokenStr := ctx.Request().Header.Get("authorization")
	idx := strings.Index(tokenStr, " ")
	if tokenStr == "" || idx < 0 {
		return errors.New("invalid jwt token")
	}

	tokenStr = tokenStr[idx+1:]
	if err := s.AuthUtil.ValidateJWTToken(tokenStr); err != nil {
		return err
	}

	userId, err := s.AuthUtil.GetUserId(tokenStr)
	if err != nil {
		return err
	}

	// TODO: validate phone number and name

	req := generated.UpdateUserProfileJSONRequestBody{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Success: false,
			Message: "Invalid Input.",
		})
	}

	if err = s.UserUsecase.UpdateUserProfile(ctx.Request().Context(), userId, req); err != nil {
		return err
	}

	resp := generated.UpdateUserProfileResponse{
		Success: true,
		Message: "successfully update user profile",
	}

	return ctx.JSON(http.StatusOK, resp)
}
