package clases

// Task es una estructura que representa una tarea con un id, input y estado.
type Tarea struct {
	IdTask int    `json:"idTask"`
	Input  string `json:"input"`
	Status int    `json:"status"` // 100: pending, 200: executing y 300: finished
}
