package classes

type PostTaskBody struct {
	Url        string   `json:"url"`
	Parameters []string `json:"parameters"`
}

// esto lo vamos a usar para una llamada post que obtenga varios taskId
type GetResultBody struct {
	TaskId []string `json:"taskId"`
}
