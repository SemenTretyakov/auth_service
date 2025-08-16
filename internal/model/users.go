package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64
	Name      string
	Email     string
	Role      int8
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserFields struct {
	Name            string
	Email           string
	Role            int32
	Password        string
	PasswordConfirm string
}
