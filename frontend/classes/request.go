package classes

type PostTaskBody struct {
	Url        string   `json:"url"`
	Parameters []string `json:"parameters"`
}

type GetResultBody struct {
	TaskId string `json:"taskId"`
}
