package entities

import (
	"errors"
)

var (
	ErrUnknownTaskID = errors.New("task not found")
)
