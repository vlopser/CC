package classes

import "time"

// Result es una estructura que representa un resultado con un id, output y tiempo.
type Result struct {
	IdTask int       `json:"idTask"`
	Output string    `json:"output"`
	Times  time.Time `json:"times"`
}

// todo esto para que sirve?
// NewResult es una función constructora que crea un nuevo resultado con los valores proporcionados.
func NewResult(IdTask int, Output string, Times time.Time) Result {
	return Result{
		IdTask: IdTask,
		Output: Output,
		Times:  Times,
	}
}
