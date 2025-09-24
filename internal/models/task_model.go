package models

type Task struct {
	Id         string `json:"id"`
	Payload    string `json:"payload"`
	MaxRetries int    `json:"max_retries"`
}

func ToTask(id, payload string, maxRetries int) *Task {
	return &Task{
		Id:         id,
		Payload:    payload,
		MaxRetries: maxRetries,
	}
}
