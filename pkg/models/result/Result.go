package result

import (
	"time"

	"github.com/google/uuid"
)

type Result struct {
	TaskId      uuid.UUID
	Output      string
	TimeElapsed time.Duration
}
