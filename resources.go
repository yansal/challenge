package main

import "time"

type Resource struct {
	ID int `json:"id"`
}

type UserResource struct {
	Resource
	Username string `json:"username"`
}

type TaskResource struct {
	Resource
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at" json:"updated_at"`
	User        UserResource `json:"user"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Progression int
}

type CommentResource struct {
	Resource
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	User      UserResource `json:"user"`
	TaskID    int          `db:"task_id" json:"task_id"`
	Content   string       `json:"content"`
}
