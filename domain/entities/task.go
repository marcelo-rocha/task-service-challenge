package entities

import (
	"time"
)

type Task struct {
	Id           int64
	Name         string
	Summary      string
	CreationDate time.Time
	FinishDate   time.Time
	UserId       int64
}
