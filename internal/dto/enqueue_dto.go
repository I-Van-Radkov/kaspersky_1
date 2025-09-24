package dto

import (
	"encoding/json"
	"io"
)

type EnqueueRequest struct {
	Id         string `json:"id"`
	Payload    string `json:"payload"`
	MaxRetries int    `json:"max_retries"`
}

func ToEnqueueRequest(fieldBody io.ReadCloser) (*EnqueueRequest, error) {
	var req EnqueueRequest
	err := json.NewDecoder(fieldBody).Decode(&req)

	return &req, err
}
