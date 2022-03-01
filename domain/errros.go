package domain

import "errors"

var ErrInvalidID = errors.New("invalid Id")
var ErrTaskNotFound = errors.New("task not found")
var ErrTaskAlreadyFinalized = errors.New("task already finalized")
