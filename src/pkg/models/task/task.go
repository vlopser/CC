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

func (s Status) String() string {
	switch s {
	case PENDING:
		return "PENDING"
	case EXECUTING:
		return "EXECUTING"
	case FINISHED:
		return "FINISHED"
	case FINISHED_ERRORS:
		return "FINISHED_WITH_ERRORS"
	}

	return ""
}

type Task struct {
	TaskId     uuid.UUID `json:"TaskId"`
	UserId     string
	RepoUrl    string   `json:"repoUrl"`
	Parameters []string `json:"parameters"`
	Input      string   //cambiar a RepoUrl
	Status     Status   `json:"status"`
}

const (
	REPO_DIR    = "/repo/"
	RESULT_DIR  = "/result/"
	STDOUT_FILE = "stdout.txt"
	STDERR_FILE = "stderr.txt"
	ERROR_FILE  = "error.txt"
)
