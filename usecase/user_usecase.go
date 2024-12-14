package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/labstack/gommon/log"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	UserRepository repository.UserRepositoryInterface
}

type UserUsecaseOptions struct {
	UserRepository repository.UserRepositoryInterface
}

func NewUserUsecase(opts UserUsecaseOptions) *UserUsecase {
	u := &UserUsecase{UserRepository: opts.UserRepository}
	return u
}

func (u UserUsecase) CreateUser(ctx context.Context, payload generated.RegisterUserJSONRequestBody) (model.User, error) {
	existingUser, err := u.UserRepository.GetUserByPhoneNumber(ctx, payload.PhoneNumber)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return model.User{}, utils.WrapWithCode(err, utils.ErrorCode(http.StatusInternalServerError), "")
	}

	if existingUser.Id > 0 {
		err = fmt.Errorf("user with phone number %s already exist", payload.PhoneNumber)
		log.Error(err)
		return model.User{}, utils.WrapWithCode(err, utils.ErrorCode(http.StatusConflict), "")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return model.User{}, utils.WrapWithCode(err, utils.ErrorCode(http.StatusInternalServerError), "")
	}

	payload.Password = string(hashedPassword)
	userId, err := u.UserRepository.CreateUser(ctx, payload)
	if err != nil {
		log.Error(err)
		return model.User{}, utils.WrapWithCode(err, utils.ErrorCode(http.StatusInternalServerError), "")
	}

	user := model.User{
		Id:          userId,
		FullName:    payload.FullName,
		PhoneNumber: payload.PhoneNumber,
	}

	return user, nil
}

func (u UserUsecase) GetUserProfile(ctx context.Context, userId int64) (model.User, error) {
	user, err := u.UserRepository.GetUserById(ctx, userId)
	if err != nil {
		log.Error(err)
		return model.User{}, utils.WrapWithCode(err, utils.ErrorCode(http.StatusForbidden), "")
	}

	return user, nil
}

func (u UserUsecase) UpdateUserProfile(ctx context.Context, userId int64, payload generated.UpdateUserProfileJSONRequestBody) error {
	user, err := u.UserRepository.GetUserById(ctx, userId)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return utils.WrapWithCode(err, utils.ErrorCode(http.StatusNotFound), "")
		}

		return utils.WrapWithCode(err, utils.ErrorCode(http.StatusInternalServerError), "")
	}

	if payload.PhoneNumber != "" && payload.PhoneNumber != user.PhoneNumber {
		existingUser, err := u.UserRepository.GetUserByPhoneNumber(ctx, payload.PhoneNumber)
		if err != nil && err != sql.ErrNoRows {
			log.Error(err)
			return utils.WrapWithCode(err, utils.ErrorCode(http.StatusInternalServerError), "")
		}

		if existingUser.Id > 0 {
			err = fmt.Errorf("phone number is already used by another user")
			log.Error(err)
			return utils.WrapWithCode(err, utils.ErrorCode(http.StatusConflict), "")
		}
	}

	if err := u.UserRepository.UpdateUserProfile(ctx, userId, payload); err != nil {
		return utils.WrapWithCode(err, utils.ErrorCode(http.StatusInternalServerError), "")
	}

	return nil
}
