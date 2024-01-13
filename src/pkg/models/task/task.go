package task

import (
	"github.com/google/uuid"
)

type Status int

// Task es una estructura que representa una tarea con un id, input y estado.
type Task struct {
	TaskId     uuid.UUID `json:"TaskId"`
	UserMail   string
	RepoUrl    string   `json:"repoUrl"`
	Parameters []string `json:"parameters"`
	Status     Status   `json:"status"` // 100: pending, 200: executing y 300: finished
}
