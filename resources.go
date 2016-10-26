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
	Progression int          `json:"progression"`
}

type TasksByCreatedAt []TaskResource

func (t TasksByCreatedAt) Len() int           { return len(t) }
func (t TasksByCreatedAt) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TasksByCreatedAt) Less(i, j int) bool { return t[i].CreatedAt.After(t[j].CreatedAt) }

type CommentResource struct {
	Resource
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	User      UserResource `json:"user"`
	TaskID    int          `db:"task_id" json:"task_id"`
	Content   string       `json:"content"`
}

type CommentsByCreatedAt []CommentResource

func (c CommentsByCreatedAt) Len() int           { return len(c) }
func (c CommentsByCreatedAt) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c CommentsByCreatedAt) Less(i, j int) bool { return c[i].CreatedAt.After(c[j].CreatedAt) }
