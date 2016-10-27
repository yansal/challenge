package main

import "time"

// Model is a base model with fields common to all models
type Model struct {
	ID        int
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// User is a model that represents a user in the database
type User struct {
	Model
	Username string
	Token    string
}

// Task is a model that represents a task in the database
type Task struct {
	Model
	UserID      int    `json:",omitempty"`
	Name        string `binding:"required"`
	Description string
	Progression int
}

// Comment is a model that represents a comment in the database
type Comment struct {
	Model
	TaskID  int    `json:",omitempty"`
	UserID  int    `json:",omitempty"`
	Content string `binding:"required"`
}
