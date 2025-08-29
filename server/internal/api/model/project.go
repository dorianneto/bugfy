package model

import "time"

type RequestCreateProject struct {
	Title string `json:"title"`
}

type ResponseCreateProject struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
