// This file contains the repository implementation layer.
package repository

import (
	"context"
	"database/sql"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/labstack/gommon/log"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type UserRepository struct {
	Db *sql.DB
}

type UserRepositoryOptions struct {
	DB *sql.DB
}

func NewUserRepository(opts UserRepositoryOptions) *UserRepository {
	return &UserRepository{Db: opts.DB}
}

func (r *UserRepository) CreateUser(ctx context.Context, payload generated.RegisterUserJSONRequestBody) (int64, error) {
	var id int64
	query := "INSERT INTO users(full_name, phone_number, password) VALUES ($1, $2, $3) RETURNING id;"
	err := r.Db.QueryRow(query, payload.FullName, payload.PhoneNumber, payload.Password).Scan(&id)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id int64) (model.User, error) {
	user := model.User{}
	err := r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number, password FROM users WHERE id = $1;", id).
		Scan(&user.Id, &user.FullName, &user.PhoneNumber, &user.Password)
	if err != nil {
		log.Error(err)
		return user, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (model.User, error) {
	user := model.User{}
	err := r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number, password FROM users WHERE phone_number = $1;", phoneNumber).
		Scan(&user.Id, &user.FullName, &user.PhoneNumber, &user.Password)
	if err != nil {
		log.Error(err)
		return user, err
	}

	return user, nil
}

func (r *UserRepository) IncrementUserLoginCount(ctx context.Context, id int64) error {
	query := "UPDATE users SET login_count = COALESCE(login_count, 0) + 1, updated_at = NOW() WHERE id = $1"
	if err := r.Db.QueryRow(query, id).Err(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (r *UserRepository) UpdateUserProfile(ctx context.Context, id int64, payload generated.UpdateUserProfileJSONRequestBody) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Update("users").Where(sq.Eq{"id": id})

	if payload.FullName != "" {
		query = query.Set("full_name", payload.FullName)
	}

	if payload.PhoneNumber != "" {
		query = query.Set("phone_number", payload.PhoneNumber)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		log.Error(err)
		return err
	}

	if err := r.Db.QueryRow(sql, args...).Err(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
