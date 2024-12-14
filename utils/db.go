package utils

import "database/sql"

type DBOptions struct {
	DSN string
}

func InitDB(opt DBOptions) (*sql.DB, error) {
	db, err := sql.Open("postgres", opt.DSN)
	if err != nil {
		return nil, err
	}

	return db, nil
}
