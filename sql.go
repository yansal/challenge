package main

const (
	create = `CREATE TABLE users(
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
	drop          = `DROP TABLE comments; DROP TABLE tasks; DROP TABLE users;`
	insertUser    = `INSERT INTO users ("username", "token") VALUES ($1, $2);`
	insertTask    = `INSERT INTO tasks ("name", "user_id", "description") VALUES ($1, $2, $3);`
	insertComment = `INSERT INTO comments ("user_id", "task_id", "content") VALUES ($1, $2, $3)`
)
