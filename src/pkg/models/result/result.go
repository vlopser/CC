package result

import (
	"time"

	"github.com/google/uuid"
)

type Result struct {
	TaskId      uuid.UUID
	Files       []string
	TimeElapsed time.Duration
}
