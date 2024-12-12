package usecase

import (
	"context"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"
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
	userId, err := u.Repository.CreateUser(ctx, payload)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		Id:          userId,
		FullName:    payload.FullName,
		PhoneNumber: payload.PhoneNumber,
	}

	return user, nil
}
