package entities

import (
	"time"
)

const (
	MaxSummarySize = 2500
)

type Task struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	Summary      string    `json:"summary"`
	CreationDate time.Time `json:"creation_date"`
	FinishDate   time.Time `json:"finish_date"`
	UserId       int64     `json:"user_id"`
}
