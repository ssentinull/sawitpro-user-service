package usecase

import (
	"context"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
)

type UsecaseInterface interface {
	CreateUser(ctx context.Context, payload generated.RegisterUserJSONRequestBody) (model.User, error)
}
