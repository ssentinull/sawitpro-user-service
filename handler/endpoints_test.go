package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/mocks"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_AuthLogin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		rec := httptest.NewRecorder()

		id := int64(1)
		phoneNumber := "+628123456782"
		password := "password"

		user := model.User{
			Id:          id,
			PhoneNumber: phoneNumber,
			Password:    password,
		}

		payload := generated.AuthLoginJSONRequestBody{
			PhoneNumber: phoneNumber,
			Password:    password,
		}

		payloadJSON, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(payloadJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		mockAuthUsecase := mocks.NewMockAuthUsecaseInterface(ctrl)
		mockAuthUsecase.EXPECT().LoginUser(gomock.Any(), payload).Times(1).Return(user, "jwt", nil)

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{AuthUsecase: mockAuthUsecase})
		s.AuthLogin(c)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	})

	t.Run("failed - missing required fields", func(t *testing.T) {
		e := echo.New()
		rec := httptest.NewRecorder()

		form := make(url.Values)
		form.Add("phone_numbers", "+6285912345678")

		payload := strings.NewReader(form.Encode())

		req := httptest.NewRequest(http.MethodPost, "/auth/login", payload)
		req.Header.Set("content-type", "application/x-www-form-urlencoded")

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{})
		s.AuthLogin(c)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	})

	t.Run("failed - invalid fields", func(t *testing.T) {
		e := echo.New()
		rec := httptest.NewRecorder()

		form := make(url.Values)
		form.Add("phone_number", "085912345678")
		form.Add("password", "password")

		payload := strings.NewReader(form.Encode())

		req := httptest.NewRequest(http.MethodPost, "/auth/login", payload)
		req.Header.Set("content-type", "application/x-www-form-urlencoded")

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{})
		s.AuthLogin(c)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	})

	t.Run("failed - login user return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		rec := httptest.NewRecorder()

		payload := generated.AuthLoginJSONRequestBody{
			PhoneNumber: "+628123456782",
			Password:    "password",
		}

		payloadJSON, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(payloadJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		mockAuthUsecase := mocks.NewMockAuthUsecaseInterface(ctrl)
		mockAuthUsecase.EXPECT().LoginUser(gomock.Any(), gomock.Any()).
			Times(1).Return(model.User{}, "", utils.NewErrorWithCode(http.StatusInternalServerError, "usecase error"))

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{AuthUsecase: mockAuthUsecase})
		s.AuthLogin(c)

		require.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
	})
}

func TestHandler_RegisterUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		rec := httptest.NewRecorder()

		id := int64(1)
		fullName := "John Doe"
		phoneNumber := "+628123456782"
		password := "passworD!1"

		user := model.User{
			Id:          id,
			PhoneNumber: phoneNumber,
			Password:    password,
		}

		payload := generated.RegisterUserJSONRequestBody{
			FullName:    fullName,
			PhoneNumber: phoneNumber,
			Password:    password,
		}

		payloadJSON, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(payloadJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		mockUserUsecase := mocks.NewMockUserUsecaseInterface(ctrl)
		mockUserUsecase.EXPECT().CreateUser(gomock.Any(), payload).Times(1).Return(user, nil)

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{UserUsecase: mockUserUsecase})
		s.RegisterUser(c)

		require.Equal(t, http.StatusCreated, rec.Result().StatusCode)
	})

	t.Run("failed - missing required fields", func(t *testing.T) {
		e := echo.New()
		rec := httptest.NewRecorder()

		payload := generated.RegisterUserJSONRequestBody{FullName: "fullName"}
		payloadJSON, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(payloadJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{})
		s.AuthLogin(c)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	})

	t.Run("failed - invalid fields", func(t *testing.T) {
		e := echo.New()
		rec := httptest.NewRecorder()

		payload := generated.RegisterUserJSONRequestBody{
			FullName:    "fullName",
			PhoneNumber: "phoneNumber",
			Password:    "password",
		}

		payloadJSON, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(payloadJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{})
		s.AuthLogin(c)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	})

	t.Run("failed - create user return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		rec := httptest.NewRecorder()

		fullName := "John Doe"
		phoneNumber := "+628123456782"
		password := "passworD!1"

		payload := generated.RegisterUserJSONRequestBody{
			FullName:    fullName,
			PhoneNumber: phoneNumber,
			Password:    password,
		}

		payloadJSON, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(payloadJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		mockUserUsecase := mocks.NewMockUserUsecaseInterface(ctrl)
		mockUserUsecase.EXPECT().CreateUser(gomock.Any(), payload).
			Times(1).Return(model.User{}, utils.NewErrorWithCode(http.StatusInternalServerError, "usecase error"))

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{UserUsecase: mockUserUsecase})
		s.RegisterUser(c)

		require.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
	})
}
