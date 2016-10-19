package main

import "time"

type UserResource struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type TaskResource struct {
	ID          int          `json:"id"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	User        UserResource `json:"user"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
}

type CommentResource struct {
	ID        int          `json:"id"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	User      UserResource `json:"user"`
	TaskID    int          `db:"task_id" json:"task_id"`
	Content   string       `json:"content"`
}
