package domain

import "time"

type JobRecord struct {
	Url string
}

type Job struct {
	Id         int
	UserId     int64
	Url        string
	Status     string
	Platform   string
	FileSize   int64
	ErrorMsg   string
	CreatedAt  time.Time
	FinishedAt time.Time
}
