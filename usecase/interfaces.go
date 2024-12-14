package usecase

import (
	"context"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
)

type AuthUsecaseInterface interface {
	LoginUser(ctx context.Context, payload generated.AuthLoginJSONRequestBody) (model.User, string, error)
}

type UserUsecaseInterface interface {
	CreateUser(ctx context.Context, payload generated.RegisterUserJSONRequestBody) (model.User, error)
	GetUserProfile(ctx context.Context, userId int64) (model.User, error)
}
