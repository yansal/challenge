package main

const (
	create = `CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = CURRENT_TIMESTAMP;
	RETURN NEW;
END;
$$ language 'plpgsql';

	CREATE TABLE users(
	"id" SERIAL PRIMARY KEY,
	"created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"username" TEXT,
	"token" TEXT
);

CREATE TRIGGER update_users_updated_at BEFORE UPDATE
ON users FOR EACH ROW EXECUTE PROCEDURE
update_updated_at_column();

CREATE TABLE tasks(
	"id" SERIAL PRIMARY KEY,
	"created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"name" TEXT NOT NULL,
	"user_id" SERIAL REFERENCES users(id),
	"description" TEXT
);

CREATE TRIGGER update_tasks_updated_at BEFORE UPDATE
ON tasks FOR EACH ROW EXECUTE PROCEDURE
update_updated_at_column();

CREATE TABLE comments(
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
	insertTask    = `INSERT INTO tasks ("name", "user_id", "description") VALUES ($1, $2, $3);`
	insertComment = `INSERT INTO comments ("user_id", "task_id", "content") VALUES ($1, $2, $3)`
)
