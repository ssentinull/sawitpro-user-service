package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"

	_ "github.com/lib/pq"
)

func (r *Repository) CreateUser(ctx context.Context, payload generated.RegisterUserJSONRequestBody) (int64, error) {
	var id int64
	query := "INSERT INTO users(full_name, phone_number, password) VALUES ($1, $2, $3) RETURNING id;"
	err := r.Db.QueryRow(query, payload.FullName, payload.PhoneNumber, payload.Password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) GetUserById(ctx context.Context, id int64) (model.User, error) {
	user := model.User{}
	err := r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number, password FROM users WHERE id = $1;", id).
		Scan(&user.Id, &user.FullName, &user.PhoneNumber, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (model.User, error) {
	user := model.User{}
	err := r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number, password FROM users WHERE phone_number = $1;", phoneNumber).
		Scan(&user.Id, &user.FullName, &user.PhoneNumber, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *Repository) IncrementUserLoginCount(ctx context.Context, id int64) error {
	query := "UPDATE users SET login_count = COALESCE(login_count, 0) + 1, updated_at = NOW() WHERE id = $1"
	if err := r.Db.QueryRow(query, id).Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateUserProfile(ctx context.Context, id int64, payload generated.UpdateUserProfileJSONRequestBody) error {

	switch {
	case payload.FullName != "" && payload.PhoneNumber != "":
		query := "UPDATE users SET full_name = $1, phone_number = $2 WHERE id = $3;"
		if err := r.Db.QueryRow(query, payload.FullName, payload.PhoneNumber, id).Err(); err != nil {
			return err
		}

	case payload.FullName != "":
		query := "UPDATE users SET full_name = $1 WHERE id = $2;"
		if err := r.Db.QueryRow(query, payload.FullName, id).Err(); err != nil {
			return err
		}

	case payload.PhoneNumber != "":
		query := "UPDATE users SET phone_number = $1 WHERE id = $2;"
		if err := r.Db.QueryRow(query, payload.PhoneNumber, id).Err(); err != nil {
			return err
		}
	}

	return nil
}
