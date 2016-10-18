package main

import "time"

type Model struct {
	ID        int64     `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type User struct {
	Model
	Username string
	Token    string
}

type Task struct {
	Model
	UserID      int64  `json:",omitempty"`
	Name        string `json:",omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
}

type ByCreatedAt []Task

func (t ByCreatedAt) Len() int           { return len(t) }
func (t ByCreatedAt) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByCreatedAt) Less(i, j int) bool { return t[i].CreatedAt.Before(t[j].CreatedAt) }
