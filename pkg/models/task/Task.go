package task

import (
	"github.com/google/uuid"
)

type Status int

const (
	PENDING   Status = 100
	EXECUTING Status = 200
	FINISHED  Status = 300
)

type Task struct {
	TaskId uuid.UUID
	Input  string
	Status Status
}
