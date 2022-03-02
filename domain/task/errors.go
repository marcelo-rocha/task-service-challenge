package task

import (
	"errors"
)

var (
	ErrNotAllowed = errors.New("the user is not allowed to see this task")
)
