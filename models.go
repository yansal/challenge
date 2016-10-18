package main

import "time"

type Model struct {
	ID        int64      `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type User struct {
	Model
	Username string
	Token    string
}

type Task struct {
	Model
	UserID      int64
	Name        string `json:"name,omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
}
