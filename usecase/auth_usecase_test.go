package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/mocks"
	mockUtils "github.com/SawitProRecruitment/UserService/mocks/utils"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
)

func TestAuthUsecase_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockAuthUtil := mockUtils.NewMockAuthInterface(ctrl)
	mockCryptUtil := mockUtils.NewMockCryptInterface(ctrl)

	authUsecase := NewAuthUsecase(AuthUsecaseOptions{
		UserRepository: mockUserRepo,
		AuthUtil:       mockAuthUtil,
		CryptUtil:      mockCryptUtil,
	})

	id := int64(1)
	phoneNumber := "+6285912345678"
	password := "password"
	jwtToken := "thisisjwt"

	payload := generated.AuthLoginRequest{
		PhoneNumber: phoneNumber,
		Password:    password,
	}

	user := model.User{
		Id:          id,
		PhoneNumber: phoneNumber,
		Password:    password,
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByPhoneNumber(ctx, phoneNumber).Times(1).Return(user, nil)
		mockCryptUtil.EXPECT().CompareHashAndPassword([]byte(password), []byte(password)).Times(1).Return(nil)
		mockAuthUtil.EXPECT().GenerateJWTToken(user).Times(1).Return(jwtToken, nil)
		mockUserRepo.EXPECT().IncrementUserLoginCount(ctx, id).Times(1).Return(errors.New("db error"))

		resUser, resToken, err := authUsecase.LoginUser(ctx, payload)
		require.NoError(t, err)
		require.NotEmpty(t, resUser)
		require.Equal(t, jwtToken, resToken)
	})

	t.Run("failed - get user by phone number return error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByPhoneNumber(ctx, phoneNumber).Times(1).Return(model.User{}, errors.New("repo error"))

		resUser, resToken, err := authUsecase.LoginUser(ctx, payload)
		require.Error(t, err)
		require.Empty(t, resUser)
		require.Zero(t, resToken)
	})

	t.Run("failed - user password doesnt match", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByPhoneNumber(ctx, phoneNumber).Times(1).Return(user, nil)
		mockCryptUtil.EXPECT().CompareHashAndPassword([]byte(password), []byte(password)).Times(1).Return(errors.New("password doesnt match"))

		resUser, resToken, err := authUsecase.LoginUser(ctx, payload)
		require.Error(t, err)
		require.Empty(t, resUser)
		require.Zero(t, resToken)
	})

	t.Run("failed - failed generate jwt", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByPhoneNumber(ctx, phoneNumber).Times(1).Return(user, nil)
		mockCryptUtil.EXPECT().CompareHashAndPassword([]byte(password), []byte(password)).Times(1).Return(nil)
		mockAuthUtil.EXPECT().GenerateJWTToken(user).Times(1).Return("", errors.New("failed generate jwt token"))

		resUser, resToken, err := authUsecase.LoginUser(ctx, payload)
		require.Error(t, err)
		require.Empty(t, resUser)
		require.Zero(t, resToken)
	})
}
