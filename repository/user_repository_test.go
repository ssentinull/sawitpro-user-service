package repository

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	ctx := context.TODO()
	userRepo := NewUserRepository(UserRepositoryOptions{DB: db})

	id := int64(10)
	fullName := "John Doe"
	password := "password"
	phoneNumber := "+6285912345678"

	payload := generated.RegisterUserRequest{
		FullName:    fullName,
		Password:    password,
		PhoneNumber: phoneNumber,
	}

	query := "INSERT INTO users(full_name, phone_number, password) VALUES ($1, $2, $3) RETURNING id;"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(strconv.FormatInt(id, 10))
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(fullName, phoneNumber, password).WillReturnRows(rows)

		resId, err := userRepo.CreateUser(ctx, payload)
		require.NoError(t, err)
		require.Equal(t, id, resId)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(fullName, phoneNumber, password).WillReturnError(errors.New("db error"))

		resId, err := userRepo.CreateUser(ctx, payload)
		require.Error(t, err)
		require.Zero(t, resId)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})
}

func TestUserRepository_GetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	ctx := context.TODO()
	userRepo := NewUserRepository(UserRepositoryOptions{DB: db})

	id := int64(10)
	fullName := "John Doe"
	password := "password"
	phoneNumber := "+6285912345678"

	query := "SELECT id, full_name, phone_number, password FROM users WHERE id = $1;"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "full_name", "phone_number", "password"}).
			AddRow(strconv.FormatInt(id, 10), fullName, phoneNumber, password)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(rows)

		resUser, err := userRepo.GetUserById(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, resUser.Id)
		require.Equal(t, fullName, resUser.FullName)
		require.Equal(t, phoneNumber, resUser.PhoneNumber)
		require.Equal(t, password, resUser.Password)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errors.New("db error"))

		resUser, err := userRepo.GetUserById(ctx, id)
		require.Error(t, err)
		require.Empty(t, resUser)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})
}

func TestUserRepository_GetUserByPhoneNumber(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	ctx := context.TODO()
	userRepo := NewUserRepository(UserRepositoryOptions{DB: db})

	id := int64(10)
	fullName := "John Doe"
	password := "password"
	phoneNumber := "+6285912345678"

	query := "SELECT id, full_name, phone_number, password FROM users WHERE phone_number = $1;"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "full_name", "phone_number", "password"}).
			AddRow(strconv.FormatInt(id, 10), fullName, phoneNumber, password)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(phoneNumber).WillReturnRows(rows)

		resUser, err := userRepo.GetUserByPhoneNumber(ctx, phoneNumber)
		require.NoError(t, err)
		require.Equal(t, id, resUser.Id)
		require.Equal(t, fullName, resUser.FullName)
		require.Equal(t, phoneNumber, resUser.PhoneNumber)
		require.Equal(t, password, resUser.Password)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(phoneNumber).WillReturnError(errors.New("db error"))

		resUser, err := userRepo.GetUserByPhoneNumber(ctx, phoneNumber)
		require.Error(t, err)
		require.Empty(t, resUser)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})
}

func TestUserRepository_IncrementUserLoginCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	ctx := context.TODO()
	userRepo := NewUserRepository(UserRepositoryOptions{DB: db})

	id := int64(10)
	query := "UPDATE users SET login_count = COALESCE(login_count, 0) + 1, updated_at = NOW() WHERE id = $1"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(strconv.FormatInt(id, 10))
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(rows)

		err := userRepo.IncrementUserLoginCount(ctx, id)
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("failed", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errors.New("db error"))

		err := userRepo.IncrementUserLoginCount(ctx, id)
		require.Error(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})
}

func TestUserRepository_UpdateUserProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	ctx := context.TODO()
	userRepo := NewUserRepository(UserRepositoryOptions{DB: db})

	id := int64(10)
	fullName := "John Doe"
	phoneNumber := "+6285912345678"

	t.Run("success", func(t *testing.T) {
		payload := generated.UpdateUserProfileJSONRequestBody{
			FullName:    fullName,
			PhoneNumber: phoneNumber,
		}

		query := "UPDATE users SET full_name = $1, phone_number = $2 WHERE id = $3"

		rows := sqlmock.NewRows([]string{"id"}).AddRow(strconv.FormatInt(id, 10))
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(fullName, phoneNumber, id).WillReturnRows(rows)

		err := userRepo.UpdateUserProfile(ctx, id, payload)
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("failed", func(t *testing.T) {
		payload := generated.UpdateUserProfileJSONRequestBody{PhoneNumber: phoneNumber}
		query := "UPDATE users SET phone_number = $1 WHERE id = $2"

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(phoneNumber, id).WillReturnError(errors.New("db error"))

		err := userRepo.UpdateUserProfile(ctx, id, payload)
		require.Error(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})
}
