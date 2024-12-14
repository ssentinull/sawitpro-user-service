// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
)

type RepositoryInterface interface {
	GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error)
	CreateUser(ctx context.Context, payload generated.RegisterUserJSONRequestBody) (int64, error)
	GetUserById(ctx context.Context, id int64) (model.User, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (model.User, error)
	IncrementUserLoginCount(ctx context.Context, id int64) error
}
