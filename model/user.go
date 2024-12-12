package model

import (
	"time"

	"github.com/guregu/null/v5"
)

type User struct {
	Id          int64
	FullName    string
	PhoneNumber string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   null.Time
	DeletedAt   null.Time
}
