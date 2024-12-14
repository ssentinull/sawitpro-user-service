package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	Repository repository.RepositoryInterface
}

type UserUsecaseOptions struct {
	Repository repository.RepositoryInterface
}

func NewUserUsecase(opts UserUsecaseOptions) *UserUsecase {
	u := &UserUsecase{Repository: opts.Repository}
	return u
}

func (u UserUsecase) CreateUser(ctx context.Context, payload generated.RegisterUserJSONRequestBody) (model.User, error) {
	// TODO: check if email is duplicate

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	payload.Password = string(hashedPassword)

	userId, err := u.Repository.CreateUser(ctx, payload)
	if err != nil {
		// TODO: implement stacktrace
		return model.User{}, err
	}

	user := model.User{
		Id:          userId,
		FullName:    payload.FullName,
		PhoneNumber: payload.PhoneNumber,
	}

	return user, nil
}

func (u UserUsecase) GetUserProfile(ctx context.Context, userId int64) (model.User, error) {
	user, err := u.Repository.GetUserById(ctx, userId)
	if err != nil {
		// TODO: implement stacktrace
		return model.User{}, err
	}

	return user, nil
}

func (u UserUsecase) UpdateUserProfile(ctx context.Context, userId int64, payload generated.UpdateUserProfileJSONRequestBody) error {
	user, err := u.Repository.GetUserById(ctx, userId)
	if err != nil {
		return err
	}

	if user.Id <= 0 {
		return errors.New("user doesnt exist")
	}

	if payload.PhoneNumber != "" && payload.PhoneNumber != user.PhoneNumber {
		userByPhoneNumber, err := u.Repository.GetUserByPhoneNumber(ctx, payload.PhoneNumber)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if userByPhoneNumber.Id > 0 {
			return errors.New("phone number is used by another user")
		}
	}

	if err := u.Repository.UpdateUserProfile(ctx, userId, payload); err != nil {
		return err
	}

	return nil
}
