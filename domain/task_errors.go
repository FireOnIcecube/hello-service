package domain

import "errors"

var ErrTaskNotFound = errors.New("task not found")

var ErrTaskTitleRequired = errors.New("task title is required")

var ErrTaskPositiveInteger = errors.New("task id must be Positive Integer")
