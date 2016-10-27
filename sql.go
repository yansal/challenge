package main

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func create() {
	db.MustExec(`CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = CURRENT_TIMESTAMP;
	RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE users (
	"id" SERIAL PRIMARY KEY,
	"created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"username" TEXT,
	"token" TEXT
);

CREATE TRIGGER update_users_updated_at BEFORE UPDATE
ON users FOR EACH ROW EXECUTE PROCEDURE
update_updated_at_column();

CREATE TABLE tasks (
	"id" SERIAL PRIMARY KEY,
	"created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"name" TEXT NOT NULL,
	"user_id" SERIAL REFERENCES users(id),
	"description" TEXT,
	"progression" integer DEFAULT 0
);

CREATE TRIGGER update_tasks_updated_at BEFORE UPDATE
ON tasks FOR EACH ROW EXECUTE PROCEDURE
update_updated_at_column();

CREATE TABLE comments (
	"id" SERIAL PRIMARY KEY,
	"created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"user_id" SERIAL REFERENCES users(id),
	"task_id" SERIAL REFERENCES tasks(id),
	"content" TEXT
);

CREATE TRIGGER update_comments_updated_at BEFORE UPDATE
ON comments FOR EACH ROW EXECUTE PROCEDURE
update_updated_at_column();`)
}

var (
	insertUser,
	insertTask,
	insertComment,
	selectTasks,
	selectTaskWhereID,
	selectTasksWhereUserID,
	selectCommentsWhereTaskID,
	selectCommentsWhereUserID,
	selectUsersWhereToken,
	updateTasksName,
	updateTasksDescription,
	updateTasksProgression *sqlx.Stmt
)

func prepare() {
	var err error
	insertUser, err = db.Preparex(`INSERT INTO users ("username", "token") VALUES ($1, $2);`)
	if err != nil {
		log.Fatal(err)
	}

	insertTask, err = db.Preparex(`INSERT INTO tasks ("name", "user_id", "description", "progression") VALUES ($1, $2, $3, $4);`)
	if err != nil {
		log.Fatal(err)
	}

	insertComment, err = db.Preparex(`INSERT INTO comments ("user_id", "task_id", "content") VALUES ($1, $2, $3)`)
	if err != nil {
		log.Fatal(err)
	}

	selectTasks, err = db.Preparex(`SELECT tasks.id, tasks.created_at, tasks.updated_at, tasks.name, tasks.description, tasks.progression, users.id AS "user.id", users.username AS "user.username"
		FROM tasks JOIN users ON tasks.user_id = users.id
		ORDER BY created_at DESC;`)
	if err != nil {
		log.Fatal(err)
	}

	selectTaskWhereID, err = db.Preparex(`SELECT tasks.id, tasks.created_at, tasks.updated_at, tasks.name, tasks.description, tasks.progression, users.id AS "user.id", users.username AS "user.username"
		FROM tasks JOIN users ON tasks.user_id = users.id
		WHERE tasks.id = $1;`)
	if err != nil {
		log.Fatal(err)
	}

	selectTasksWhereUserID, err = db.Preparex(`SELECT tasks.id, tasks.created_at, tasks.updated_at, tasks.name, tasks.description, tasks.progression, users.id AS "user.id", users.username AS "user.username"
		FROM tasks JOIN users ON tasks.user_id = users.id
		WHERE user_id = $1
		ORDER BY created_at DESC;`)
	if err != nil {
		log.Fatal(err)
	}

	selectCommentsWhereTaskID, err = db.Preparex(`SELECT comments.id, comments.created_at, comments.content, comments.task_id, users.id AS "user.id", users.username AS "user.username"
		FROM comments JOIN users ON comments.user_id = users.id
		WHERE task_id = $1
		ORDER BY created_at DESC;`)
	if err != nil {
		log.Fatal(err)
	}

	selectCommentsWhereUserID, err = db.Preparex(`SELECT comments.id, comments.created_at, comments.content, comments.task_id, users.id AS "user.id", users.username AS "user.username"
		FROM comments JOIN users ON comments.user_id = users.id
		WHERE user_id = $1
		ORDER BY created_at DESC;`)
	if err != nil {
		log.Fatal(err)
	}

	selectUsersWhereToken, err = db.Preparex(`SELECT * FROM users WHERE token = $1`)
	if err != nil {
		log.Fatal(err)
	}

	updateTasksName, err = db.Preparex(`UPDATE tasks SET name = $1 WHERE id = $2;`)
	if err != nil {
		log.Fatal(err)
	}

	updateTasksDescription, err = db.Preparex(`UPDATE tasks SET description = $1 WHERE id = $2;`)
	if err != nil {
		log.Fatal(err)
	}

	updateTasksProgression, err = db.Preparex(`UPDATE tasks SET progression = $1 WHERE id = $2;`)
	if err != nil {
		log.Fatal(err)
	}
}

func seed() {
	for _, user := range []User{
		{Username: "Alice", Token: "077000ac559e1ba0fe4f303b614f30da6306341f"},
		{Username: "Bob", Token: "ef2e253a2b4564ae949b053025c845552f2e99cc"},
	} {
		insertUser.MustExec(user.Username, user.Token)
	}
	for _, task := range []Task{
		{Name: "First task", UserID: 1, Description: "This is the first task"},
		{Name: "Second task", UserID: 1, Description: "This is the second task"},
		{Name: "Third task", UserID: 2, Description: "This is the third task"},
	} {
		insertTask.MustExec(task.Name, task.UserID, task.Description, task.Progression)
	}
	for _, comment := range []Comment{
		{TaskID: 1, UserID: 1, Content: "This is the first comment"},
		{TaskID: 1, UserID: 2, Content: "This is the second comment"},
		{TaskID: 2, UserID: 2, Content: "This is the third comment"},
	} {
		insertComment.MustExec(comment.UserID, comment.TaskID, comment.Content)
	}
}

func drop() {
	db.MustExec(`DROP TABLE comments; DROP TABLE tasks; DROP TABLE users;`)
}
