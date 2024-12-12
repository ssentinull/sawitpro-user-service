package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/generated"

	_ "github.com/lib/pq"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) CreateUser(ctx context.Context, payload generated.RegisterUserJSONRequestBody) (int64, error) {
	var id int64
	query := "INSERT INTO users(full_name, phone_number, password) VALUES ($1, $2, $3) RETURNING id;"
	err := r.Db.QueryRow(query, payload.FullName, payload.PhoneNumber, payload.Password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
