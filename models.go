package main

import "time"

type Model struct {
	ID        int       `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type User struct {
	Model
	Username string
	Token    string
}

type Task struct {
	Model
	UserID      int    `json:",omitempty"`
	Name        string `json:",omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
}

type TasksByCreatedAt []Task

func (t TasksByCreatedAt) Len() int           { return len(t) }
func (t TasksByCreatedAt) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TasksByCreatedAt) Less(i, j int) bool { return t[i].CreatedAt.Before(t[j].CreatedAt) }

type Comment struct {
	Model
	TaskID  int    `json:",omitempty"`
	UserID  int    `json:",omitempty"`
	Content string `json:",omitempty" binding:"required"`
}
type CommentsByCreatedAt []Comment

func (c CommentsByCreatedAt) Len() int           { return len(c) }
func (c CommentsByCreatedAt) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c CommentsByCreatedAt) Less(i, j int) bool { return c[i].CreatedAt.Before(c[j].CreatedAt) }
