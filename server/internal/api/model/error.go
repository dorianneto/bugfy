package model

import "time"

type RequestCreateError struct {
	ProjectID string            `json:"project_id"`
	Message   string            `json:"message"`
	Context   map[string]string `json:"context"`
}

type ResponseCreateError struct {
	ID          string            `json:"id"`
	ProjectID   string            `json:"project_id"`
	Message     string            `json:"message"`
	Type        string            `json:"type"`
	Fingerprint string            `json:"fingerprint"`
	Context     map[string]string `json:"context"`
	Timestamp   time.Time         `json:"timestamp"`
}
