package errors

import "errors"

var (
	ErrUserInvalid = errors.New("error: user_id is invalid")
	ErrTaskInvalid = errors.New("error: task_id is invalid")

	ErrUserNotFound = errors.New("error: user_id not found")
	ErrTaskNotFound = errors.New("error: task_id not found")
)
