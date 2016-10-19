package main

import "time"

type Model struct {
	ID        int       `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

type User struct {
	Model
	Username string
	Token    string
}

type Task struct {
	Model
	UserID      int    `json:",omitempty" db:"user_id"`
	Name        string `json:",omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
}

type Comment struct {
	Model
	TaskID  int    `json:",omitempty" db:"task_id"`
	UserID  int    `json:",omitempty" db:"user_id"`
	Content string `json:",omitempty" binding:"required"`
}
