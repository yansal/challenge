package main

import (
	"crypto/md5"
	"fmt"
	"time"
)

// Resource is a base resource with fields common to all resources
type Resource struct {
	ID int `json:"id"`
}

// UserResource is a resource that represents a user. It is always embeded in
// a TaskResource or in a CommentResource
type UserResource struct {
	Resource
	Username string `json:"username"`
}

// TaskResource is a resource that represents a task. It embeds a UserResource
type TaskResource struct {
	Resource
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at" json:"updated_at"`
	User        UserResource `json:"user"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Progression int          `json:"progression"`
}

// Etag returns the Etag for the resource
func (task *TaskResource) Etag() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprint(task))))
}

// TasksByCreatedAt is a type to sort a slice of TaskResource by creation date.
// It is only used in tests
type TasksByCreatedAt []TaskResource

func (t TasksByCreatedAt) Len() int           { return len(t) }
func (t TasksByCreatedAt) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TasksByCreatedAt) Less(i, j int) bool { return t[i].CreatedAt.After(t[j].CreatedAt) }

// CommentResource is a resource that represents a user. It embeds a
// UserResource
type CommentResource struct {
	Resource
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	User      UserResource `json:"user"`
	TaskID    int          `db:"task_id" json:"task_id"`
	Content   string       `json:"content"`
}

// CommentsByCreatedAt is a type to sort a slice of CommentResource by creation
// date. It is only used in tests
type CommentsByCreatedAt []CommentResource

func (c CommentsByCreatedAt) Len() int           { return len(c) }
func (c CommentsByCreatedAt) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c CommentsByCreatedAt) Less(i, j int) bool { return c[i].CreatedAt.After(c[j].CreatedAt) }
