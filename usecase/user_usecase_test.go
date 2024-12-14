package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/mocks"
	mockUtils "github.com/SawitProRecruitment/UserService/mocks/utils"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserUsecase_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCryptUtil := mockUtils.NewMockCryptInterface(ctrl)

	userUsecase := NewUserUsecase(UserUsecaseOptions{
		UserRepository: mockUserRepo,
		CryptUtil:      mockCryptUtil,
	})

	id := int64(1)
	fullName := "John Doe"
	phoneNumber := "+6285912345678"
	password := "password"

	payload := generated.RegisterUserJSONRequestBody{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Password:    password,
	}

	user := model.User{
		Id:          id,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Password:    password,
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByPhoneNumber(ctx, phoneNumber).Times(1).Return(model.User{}, nil)
		mockCryptUtil.EXPECT().GenerateFromPassword([]byte(password), bcrypt.DefaultCost).
			Times(1).Return([]byte(password), nil)
		mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Times(1).Return(id, nil)

		res, err := userUsecase.CreateUser(ctx, payload)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})

	t.Run("failed - get user by phone number return error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByPhoneNumber(ctx, phoneNumber).Times(1).Return(model.User{}, errors.New("repo error"))

		res, err := userUsecase.CreateUser(ctx, payload)
		require.Error(t, err)
		require.Empty(t, res)
	})

	t.Run("failed - phone number is already used", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByPhoneNumber(ctx, phoneNumber).Times(1).Return(user, nil)

		res, err := userUsecase.CreateUser(ctx, payload)
		require.Error(t, err)
		require.Empty(t, res)
	})

	t.Run("failed - generate hashed password failed", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByPhoneNumber(ctx, phoneNumber).Times(1).Return(model.User{}, nil)
		mockCryptUtil.EXPECT().GenerateFromPassword([]byte(password), bcrypt.DefaultCost).
			Times(1).Return([]byte{}, errors.New("password doesnt match"))

		res, err := userUsecase.CreateUser(ctx, payload)
		require.Error(t, err)
		require.Empty(t, res)
	})

	t.Run("failed - create user return error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByPhoneNumber(ctx, phoneNumber).Times(1).Return(model.User{}, nil)
		mockCryptUtil.EXPECT().GenerateFromPassword([]byte(password), bcrypt.DefaultCost).
			Times(1).Return([]byte(password), nil)
		mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Times(1).Return(int64(0), errors.New("db error"))

		res, err := userUsecase.CreateUser(ctx, payload)
		require.Error(t, err)
		require.Empty(t, res)
	})
}

func TestUserUsecase_GetUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCryptUtil := mockUtils.NewMockCryptInterface(ctrl)

	userUsecase := NewUserUsecase(UserUsecaseOptions{
		UserRepository: mockUserRepo,
		CryptUtil:      mockCryptUtil,
	})

	id := int64(1)
	fullName := "John Doe"
	phoneNumber := "+6285912345678"
	password := "password"

	user := model.User{
		Id:          id,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Password:    password,
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserById(ctx, id).Times(1).Return(user, nil)
		res, err := userUsecase.GetUserProfile(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})

	t.Run("failed", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserById(ctx, id).Times(1).Return(model.User{}, errors.New("db error"))
		res, err := userUsecase.GetUserProfile(ctx, id)
		require.Error(t, err)
		require.Empty(t, res)
	})
}

func TestUserUsecase_UpdateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCryptUtil := mockUtils.NewMockCryptInterface(ctrl)

	userUsecase := NewUserUsecase(UserUsecaseOptions{
		UserRepository: mockUserRepo,
		CryptUtil:      mockCryptUtil,
	})

	id := int64(1)
	fullName := "John Doe"
	phoneNumber := "+6285912345678"
	phoneNumberToUpdate := "+6285912345679"

	payload := generated.UpdateUserProfileJSONRequestBody{
		FullName:    fullName,
		PhoneNumber: phoneNumberToUpdate,
	}

	user := model.User{
		Id:          id,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}

	otherUser := model.User{
		Id:          int64(2),
		PhoneNumber: phoneNumberToUpdate,
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserById(ctx, id).Times(1).Return(user, nil)
		mockUserRepo.EXPECT().GetUserByPhoneNumber(ctx, phoneNumberToUpdate).Times(1).Return(model.User{}, nil)
		mockUserRepo.EXPECT().UpdateUserProfile(ctx, id, payload).Times(1).Return(nil)

		err := userUsecase.UpdateUserProfile(ctx, id, payload)
		require.NoError(t, err)
	})

	t.Run("failed - get user by id return error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserById(ctx, id).Times(1).Return(model.User{}, errors.New("db error"))
		err := userUsecase.UpdateUserProfile(ctx, id, payload)
		require.Error(t, err)
	})

	t.Run("failed - get user by id return no user", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserById(ctx, id).Times(1).Return(model.User{}, sql.ErrNoRows)
		err := userUsecase.UpdateUserProfile(ctx, id, payload)
		require.Error(t, err)
	})

	t.Run("failed - get user by phone number return error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserById(ctx, id).Times(1).Return(user, nil)
		mockUserRepo.EXPECT().GetUserByPhoneNumber(ctx, phoneNumberToUpdate).Times(1).Return(model.User{}, errors.New("db error"))

		err := userUsecase.UpdateUserProfile(ctx, id, payload)
		require.Error(t, err)
	})

	t.Run("failed - get user by phone number return error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserById(ctx, id).Times(1).Return(user, nil)
		mockUserRepo.EXPECT().GetUserByPhoneNumber(ctx, phoneNumberToUpdate).Times(1).Return(otherUser, nil)

		err := userUsecase.UpdateUserProfile(ctx, id, payload)
		require.Error(t, err)
	})

	t.Run("failed - update user profile return error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserById(ctx, id).Times(1).Return(user, nil)
		mockUserRepo.EXPECT().GetUserByPhoneNumber(ctx, phoneNumberToUpdate).Times(1).Return(model.User{}, nil)
		mockUserRepo.EXPECT().UpdateUserProfile(ctx, id, payload).Times(1).Return(errors.New("db error"))

		err := userUsecase.UpdateUserProfile(ctx, id, payload)
		require.Error(t, err)
	})
}
