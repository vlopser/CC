package result

import (
	"time"

	"github.com/google/uuid"
)

type Result struct {
	TaskId      uuid.UUID
	Files       []string
	Output      string
	Errors      string
	TimeElapsed time.Duration
}