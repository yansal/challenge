package main

import "time"

type Model struct {
	ID        int
	CreatedAt time.Time `db:"created_at"`
}

type User struct {
	Model
	Username string
	Token    string
}

type Task struct {
	Model
	UserID      int    `json:",omitempty"`
	Name        string `binding:"required"`
	Description string
}

type Comment struct {
	Model
	TaskID  int    `json:",omitempty"`
	UserID  int    `json:",omitempty"`
	Content string `binding:"required"`
}
