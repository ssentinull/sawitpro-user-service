package usecase

import (
	"context"
	"fmt"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	Repository repository.RepositoryInterface
	Auth       utils.AuthInterface
}

type AuthUsecaseOptions struct {
	Repository repository.RepositoryInterface
	Auth       utils.AuthInterface
}

func NewAuthUsecase(opts AuthUsecaseOptions) *AuthUsecase {
	u := &AuthUsecase{
		Repository: opts.Repository,
		Auth:       opts.Auth,
	}

	return u
}

func (u AuthUsecase) LoginUser(ctx context.Context, payload generated.AuthLoginJSONRequestBody) (model.User, string, error) {
	user, err := u.Repository.GetUserByPhoneNumber(ctx, payload.PhoneNumber)
	if err != nil {
		return model.User{}, "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return model.User{}, "", err
	}

	jwt, err := u.Auth.GenerateJWT(user)
	if err != nil {
		return model.User{}, "", err
	}

	if err = u.Repository.IncrementUserLoginCount(ctx, user.Id); err != nil {
		fmt.Println("error: ", err)
	}

	return user, jwt, nil
}
