package usecase

import (
	"context"

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
