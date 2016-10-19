package main

import "fmt"

const create = `CREATE TABLE users(
	"id" SERIAL PRIMARY KEY,
	"created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"username" TEXT,
	"token" TEXT
);

CREATE TABLE tasks(
	"id" SERIAL PRIMARY KEY,
	"created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"name" TEXT NOT NULL,
	"user_id" SERIAL REFERENCES users(id),
	"description" TEXT
);

CREATE TABLE comments(
	"id" SERIAL PRIMARY KEY,
	"created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"user_id" SERIAL REFERENCES users(id),
	"task_id" SERIAL REFERENCES tasks(id),
	"content" TEXT
);`

func insertUser(user User) error {
	if _, err := db.Exec(`INSERT INTO users ("username", "token") VALUES ($1, $2)`, user.Username, user.Token); err != nil {
		return fmt.Errorf("couldn't insert into users: %v", err)
	}
	return nil
}

func insertTask(task Task) error {
	if _, err := db.Exec(`INSERT INTO tasks ("name", "user_id", "description") VALUES ($1, $2, $3)`, task.Name, task.UserID, task.Description); err != nil {
		return fmt.Errorf("couldn't insert into tasks: %v", err)
	}
	return nil
}

func insertComment(comment Comment) error {
	if _, err := db.Exec(`INSERT INTO comments ("user_id", "task_id", "content") VALUES ($1, $2, $3)`, comment.UserID, comment.TaskID, comment.Content); err != nil {
		return fmt.Errorf("couldn't insert into comments: %v", err)
	}
	return nil
}
