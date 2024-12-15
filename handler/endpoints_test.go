package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/mocks"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	jwtSecretKey = `MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCySk/3GB3Ofjua
0klaMtQH0iQbbrvpMSgkJ0BVyO1kozLCZdpsW/kaXP0ccPzKcdPkmDTNbZakOkKH
zI0YggIQMjExwRVJQ+pryEIkBhbJ/lseknIxPjhdnQn740jJJiXbEdaOJGFYtcbp
FnTxGOzULCdjHcUlJqe41AU98RSyMGElJ8Rk659v9A6UZQaPJuTsmBrdvV+AZUHW
ev8Rp+No31ePyj+f/f/cW0WZ3lhwkN5rzat9U+xlNyIJICzkmDuJaSAWVeX4Lzr+
hNKCUL3F1hOxvxI47nIjK6+DrbNuc8HDjT1o2WWuTFHvxuLDgjfK31841+millvS
RqzXxTDFAgMBAAECggEAEb5FMgU+O3rNlMKrwSRrfIcZx88f4qeyjn0yMQyzhL9H
FxigLFuE3Bj3/u4cV6C3Yo8RO4aApkG+tZlnoKeY6/gKyaBPwtWq9+Swobl8vXJ+
VU8OupQuHjGO46NYGYid3izqyi+YWTCR+0gcWugiSEVH+txkw9CSgtmQLKdtlK6t
ckmp6DxBCb0aILm3YnulhLmCjfvePaiigd+W03VK9+M7yI8xmllcZYGvzdVBlL0k
NpqEABi/sZEMUsmwS49EMzua9AovfwfHJ6G5DCFpiysJFFFkFAUeBz4XVyrczq41
EBwMq4jSmJeMhLWpRMOurodU4pu6FtqnudKaSZXY2wKBgQDuYa/qIL3CuW31y70n
ZfskLdxKK/d+g8da2yYxyWIC8AZS9XIf+eJ1NMRhvqLtSHVGQpBt//UONbKI19RJ
6i0HlR4PQ4VK6XB0N3OowvcnNfQ0Uo7d+f9Sp35LR+b9OUpAdo++QkxJ36imFisX
4nHjTiuIgXCsaJhBqZk3xvTTzwKBgQC/d6oB5oU/8huF7WTfdeD15UiIA+2GwTji
SNCkrp1ZnMl2H83joFAfF9fj5RTpERdxNUsVyH5Ihe6eeMCBAU4K4oo2f5+8sY2b
jZAc1UrvKTMALZpIiRIyvAznrkGUqN9vlhIa/hB4INnVP9s1oY0rT7NzGbB3f2Xg
yB1Fl07TKwKBgHk45gtKkRUn1Los3EjfvGG+jIqPZzFH9CXI0dh5j0TtKFohhOKr
4TQ3HDKUjifaNAEBso6tncGXHu4ly0e3NSTo+LtMW8kngs8mr8M/Og4PitrcrNhG
3Eb88+V2cAmPi6nSYPCgqEjc2tdy6IEh30Z3Jv4ozNJv8hVaGJdbrn7TAoGAf/Qj
bBO21u4gUJc+Q0vOo+WvXB5r3RNBxY9tx7BdvWZXCBbnDAi1oqHXiBguqjbe2KwJ
2qvbIPJIbiU6WLwbgJC2VwdhI8PwY5TuSyaLZlq9F5BiO7lGrRsY8Ld2Yjec4kCD
JwDE1tL1YFrFTwkAg4JG5VO0p5c+6UIytbARYHMCgYBT77UKVL7NdV9zo9Qk2bJ6
xfNA1zEou6KdaRG00idhEGeXmUb5vmUsHxsL/hlUGlf7kZxoIZC814a1ibSTU7iT
bfvg0Qvy85i0uVFanyEWKZRYAUWabp+WtBoNaVHZKuW/kteQDuTExvp4IEecoPoa
5XHGaDwFf259UPSsiQDQOg==`
	jwtDurationStr = "1h"
	dummyJwt       = "dummyJwt"
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

		var response generated.AuthLoginResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.NotEmpty(t, response.Message)
		require.NotEmpty(t, response.Data.Id)
		require.NotEmpty(t, response.Data.Jwt)
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

		var response generated.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.False(t, response.Success)
		require.NotEmpty(t, response.Message)
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

		var response generated.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.False(t, response.Success)
		require.NotEmpty(t, response.Message)
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

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.False(t, response.Success)
		require.NotEmpty(t, response.Message)
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

		var response generated.RegisterUserResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.NotEmpty(t, response.Message)
		require.NotEmpty(t, response.Data.Id)
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

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.False(t, response.Success)
		require.NotEmpty(t, response.Message)
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

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.False(t, response.Success)
		require.NotEmpty(t, response.Message)
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

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.False(t, response.Success)
		require.NotEmpty(t, response.Message)
	})
}

func TestHandler_GetUserProfile(t *testing.T) {
	jwtDuration, err := time.ParseDuration(jwtDurationStr)
	require.NoError(t, err)

	auth, err := utils.InitAuth(utils.AuthOptions{
		JWTSecretKey:      jwtSecretKey,
		JWTExpiryDuration: jwtDuration,
	})
	require.NoError(t, err)

	id := int64(10)
	fullName := "John Doe"
	phoneNumber := "+6285912345678"

	user := model.User{
		Id:          id,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}

	jwt, err := auth.GenerateJWTToken(user)
	require.NoError(t, err)

	authUtil, err := utils.InitAuth(utils.AuthOptions{
		JWTSecretKey:      jwtSecretKey,
		JWTExpiryDuration: jwtDuration,
	})
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		rec := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", jwt))

		mockUserUsecase := mocks.NewMockUserUsecaseInterface(ctrl)
		mockUserUsecase.EXPECT().GetUserProfile(gomock.Any(), id).Times(1).Return(user, nil)

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{
			UserUsecase: mockUserUsecase,
			AuthUtil:    authUtil,
		})
		s.GetUserProfile(c)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)

		var response generated.GetUserProfileResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.NotEmpty(t, response.Message)
		require.NotEmpty(t, response.Data.Id)
		require.NotEmpty(t, response.Data.FullName)
		require.NotEmpty(t, response.Data.PhoneNumber)
	})

	t.Run("failed - missing jwt token", func(t *testing.T) {
		e := echo.New()
		rec := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
		req.Header.Set("authorization", "Bearer")

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{})
		s.GetUserProfile(c)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.False(t, response.Success)
		require.NotEmpty(t, response.Message)
	})

	t.Run("failed - invalid jwt token", func(t *testing.T) {
		e := echo.New()
		rec := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", dummyJwt))

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{AuthUtil: authUtil})
		s.GetUserProfile(c)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.False(t, response.Success)
		require.NotEmpty(t, response.Message)
	})

	t.Run("failed - get user profile return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		rec := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", jwt))

		mockUserUsecase := mocks.NewMockUserUsecaseInterface(ctrl)
		mockUserUsecase.EXPECT().GetUserProfile(gomock.Any(), id).
			Times(1).Return(model.User{}, utils.NewErrorWithCode(http.StatusInternalServerError, "usecase error"))

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{
			UserUsecase: mockUserUsecase,
			AuthUtil:    authUtil,
		})

		s.GetUserProfile(c)

		require.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.False(t, response.Success)
		require.NotEmpty(t, response.Message)
	})
}

func TestHandler_UpdateUserProfile(t *testing.T) {
	jwtDuration, err := time.ParseDuration(jwtDurationStr)
	require.NoError(t, err)

	auth, err := utils.InitAuth(utils.AuthOptions{
		JWTSecretKey:      jwtSecretKey,
		JWTExpiryDuration: jwtDuration,
	})
	require.NoError(t, err)

	id := int64(10)
	fullName := "John Doe"
	phoneNumber := "+6285912345678"

	user := model.User{
		Id:          id,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}

	jwt, err := auth.GenerateJWTToken(user)
	require.NoError(t, err)

	authUtil, err := utils.InitAuth(utils.AuthOptions{
		JWTSecretKey:      jwtSecretKey,
		JWTExpiryDuration: jwtDuration,
	})

	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		rec := httptest.NewRecorder()

		payload := generated.UpdateUserProfileJSONRequestBody{
			FullName:    fullName,
			PhoneNumber: phoneNumber,
		}

		payloadJSON, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/users/profile", bytes.NewReader(payloadJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", jwt))

		mockUserUsecase := mocks.NewMockUserUsecaseInterface(ctrl)
		mockUserUsecase.EXPECT().UpdateUserProfile(gomock.Any(), id, payload).Times(1).Return(nil)

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{
			UserUsecase: mockUserUsecase,
			AuthUtil:    authUtil,
		})

		s.UpdateUserProfile(c)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)

		var response generated.UpdateUserProfileResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.True(t, response.Success)
		require.NotEmpty(t, response.Message)
	})

	t.Run("failed - missing jwt token", func(t *testing.T) {
		e := echo.New()
		rec := httptest.NewRecorder()

		payload := generated.UpdateUserProfileJSONRequestBody{
			FullName:    fullName,
			PhoneNumber: phoneNumber,
		}

		payloadJSON, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/users/profile", bytes.NewReader(payloadJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("authorization", "Bearer")

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{})
		s.UpdateUserProfile(c)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.False(t, response.Success)
		require.NotEmpty(t, response.Message)
	})

	t.Run("failed - invalid jwt token", func(t *testing.T) {
		e := echo.New()
		rec := httptest.NewRecorder()

		payload := generated.UpdateUserProfileJSONRequestBody{
			FullName:    fullName,
			PhoneNumber: phoneNumber,
		}

		payloadJSON, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/users/profile", bytes.NewReader(payloadJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", dummyJwt))

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{AuthUtil: authUtil})
		s.UpdateUserProfile(c)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.False(t, response.Success)
		require.NotEmpty(t, response.Message)
	})

	t.Run("failed - invalid request body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		rec := httptest.NewRecorder()

		payload := generated.UpdateUserProfileJSONRequestBody{
			FullName:    fullName,
			PhoneNumber: "0825",
		}

		payloadJSON, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/users/profile", bytes.NewReader(payloadJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", jwt))

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{AuthUtil: authUtil})

		s.UpdateUserProfile(c)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.False(t, response.Success)
		require.NotEmpty(t, response.Message)
	})

	t.Run("failed - get user profile return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		rec := httptest.NewRecorder()

		payload := generated.UpdateUserProfileJSONRequestBody{
			FullName:    fullName,
			PhoneNumber: phoneNumber,
		}

		payloadJSON, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/users/profile", bytes.NewReader(payloadJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", jwt))

		mockUserUsecase := mocks.NewMockUserUsecaseInterface(ctrl)
		mockUserUsecase.EXPECT().UpdateUserProfile(gomock.Any(), id, payload).
			Times(1).Return(utils.NewErrorWithCode(http.StatusInternalServerError, "usecase error"))

		c := e.NewContext(req, rec)
		s := NewServer(NewServerOptions{
			UserUsecase: mockUserUsecase,
			AuthUtil:    authUtil,
		})

		s.UpdateUserProfile(c)

		require.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.False(t, response.Success)
		require.NotEmpty(t, response.Message)
	})
}
