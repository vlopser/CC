package classes

import "time"

// Result es una estructura que representa un resultado con un id, output y tiempo.
type Result struct {
	IdTask int       `json:"idTask"`
	Output string    `json:"output"`
	Times  time.Time `json:"times"`
}
