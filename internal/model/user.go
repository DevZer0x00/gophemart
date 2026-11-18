package model

import "time"

type User struct {
	Id           *uint64
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}
