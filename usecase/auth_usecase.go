package usecase

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"

	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	UserRepository repository.UserRepositoryInterface
	AuthUtil       utils.AuthInterface
}

type AuthUsecaseOptions struct {
	UserRepository repository.UserRepositoryInterface
	AuthUtil       utils.AuthInterface
}

func NewAuthUsecase(opts AuthUsecaseOptions) *AuthUsecase {
	u := &AuthUsecase{
		UserRepository: opts.UserRepository,
		AuthUtil:       opts.AuthUtil,
	}

	return u
}

func (u AuthUsecase) LoginUser(ctx context.Context, payload generated.AuthLoginJSONRequestBody) (model.User, string, error) {
	user, err := u.UserRepository.GetUserByPhoneNumber(ctx, payload.PhoneNumber)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return model.User{}, "", utils.WrapWithCode(err, utils.ErrorCode(http.StatusBadRequest), "")
		}
		return model.User{}, "", utils.WrapWithCode(err, utils.ErrorCode(http.StatusInternalServerError), "")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		log.Error(err)
		return model.User{}, "", utils.WrapWithCode(err, utils.ErrorCode(http.StatusUnauthorized), "")
	}

	jwt, err := u.AuthUtil.GenerateJWTToken(user)
	if err != nil {
		log.Error(err)
		return model.User{}, "", utils.WrapWithCode(err, utils.ErrorCode(http.StatusInternalServerError), "")
	}

	if err = u.UserRepository.IncrementUserLoginCount(ctx, user.Id); err != nil {
		log.Warn(err)
	}

	return user, jwt, nil
}
