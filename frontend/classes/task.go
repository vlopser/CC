package classes

// Task es una estructura que representa una tarea con un id, input y estado.
type Task struct {
	IdTask     int      `json:"idTask"`
	RepoUrl    string   `json:"repoUrl"`
	Parameters []string `json:"parameters"`
	Status     int      `json:"status"` // 100: pending, 200: executing y 300: finished
}

// NewTask es una función constructora que crea una nueva tarea con los valores proporcionados.
func NewTask(idTask int, repoUrl string, parameters []string) Task {
	return Task{
		IdTask:     idTask,
		RepoUrl:    repoUrl,
		Parameters: parameters,
		Status:     100,
	}
}
