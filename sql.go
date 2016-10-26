package main

const (
	create = `CREATE OR REPLACE FUNCTION update_updated_at_column()
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
update_updated_at_column();`
	drop          = `DROP TABLE comments; DROP TABLE tasks; DROP TABLE users;`
	insertUser    = `INSERT INTO users ("username", "token") VALUES ($1, $2);`
	insertTask    = `INSERT INTO tasks ("name", "user_id", "description", "progression") VALUES ($1, $2, $3, $4);`
	insertComment = `INSERT INTO comments ("user_id", "task_id", "content") VALUES ($1, $2, $3)`
)

func seed() {
	db.MustExec(create)
	for _, user := range []User{
		{Username: "Alice", Token: "077000ac559e1ba0fe4f303b614f30da6306341f"},
		{Username: "Bob", Token: "ef2e253a2b4564ae949b053025c845552f2e99cc"},
	} {
		db.MustExec(insertUser, user.Username, user.Token)
	}
	for _, task := range []Task{
		{Name: "First task", UserID: 1, Description: "This is the first task"},
		{Name: "Second task", UserID: 1, Description: "This is the second task"},
		{Name: "Third task", UserID: 2, Description: "This is the third task"},
	} {
		db.MustExec(insertTask, task.Name, task.UserID, task.Description, task.Progression)
	}
	for _, comment := range []Comment{
		{TaskID: 1, UserID: 1, Content: "This is the first comment"},
		{TaskID: 1, UserID: 2, Content: "This is the second comment"},
		{TaskID: 2, UserID: 2, Content: "This is the third comment"},
	} {
		db.MustExec(insertComment, comment.UserID, comment.TaskID, comment.Content)
	}
}
