package domain

import "time"

type User struct {
	ID        int64
	Username  string
	FirstName string
	LastName  string
	IsAdmin   bool
	IsBanned  bool
	CreatedAt time.Time
}
