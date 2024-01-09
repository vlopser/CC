package classes

type RequestBody struct {
	Url        string   `json:"url"`
	Parameters []string `json:"parameters"`
}
