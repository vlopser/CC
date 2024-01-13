package task

import (
	"github.com/google/uuid"
)

type Status int

const (
	PENDING         Status = 100
	EXECUTING       Status = 200
	FINISHED        Status = 300
	FINISHED_ERRORS Status = 400
)

type Task struct {
	TaskId uuid.UUID
	Input  string
	Status Status
}

const (
	REPO_DIR    = "/repo/"
	RESULT_DIR  = "/result/"
	STDOUT_FILE = "stdout.txt"
	STDERR_FILE = "stderr.txt"
	ERROR_FILE  = "error.txt"
)
