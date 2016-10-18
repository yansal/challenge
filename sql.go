package main

import "fmt"

func createTableUsers() error {
	if _, err := db.Exec(`CREATE TABLE users(
		"id" SERIAL PRIMARY KEY,
		"created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		"username" TEXT,
		"token" TEXT
);`); err != nil {
		return fmt.Errorf(`couldn't create "users" table: %v`, err)
	}
	return nil
}

func insertUser(user User) error {
	if _, err := db.Exec(`INSERT INTO users ("username", "token") VALUES ($1, $2)`, user.Username, user.Token); err != nil {
		return fmt.Errorf("couldn't insert into users: %v", err)
	}
	return nil
}

func selectUser(token string) (User, error) {
	row := db.QueryRow(`SELECT * FROM users WHERE token=$1`, token)
	var user User
	err := row.Scan(&user.ID, &user.CreatedAt, &user.Username, &user.Token)
	return user, err
}

func createTableTasks() error {
	if _, err := db.Exec(`CREATE TABLE tasks(
		"id" SERIAL PRIMARY KEY,
		"created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		"name" TEXT NOT NULL,
		"user_id" SERIAL REFERENCES users(id),
		"description" TEXT
);`); err != nil {
		return fmt.Errorf(`couldn't create "tasks" table: %v`, err)
	}
	return nil
}

func insertTask(task Task) error {
	if _, err := db.Exec(`INSERT INTO tasks ("name", "user_id", "description") VALUES ($1, $2, $3)`, task.Name, task.UserID, task.Description); err != nil {
		return fmt.Errorf("couldn't insert into tasks: %v", err)
	}
	return nil
}

func selectTasks() ([]Task, error) {
	rows, err := db.Query(`SELECT * FROM tasks ORDER BY "created_at";`)
	if err != nil {
		return nil, fmt.Errorf("couldn't select from tasks table: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.CreatedAt, &task.Name, &task.UserID, &task.Description); err != nil {
			return tasks, fmt.Errorf(`couldn't scan row: %v`, err)
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func selectTask(id int) (Task, error) {
	row := db.QueryRow(`SELECT * FROM tasks WHERE id=$1`, id)
	var task Task
	err := row.Scan(&task.ID, &task.CreatedAt, &task.Name, &task.UserID, &task.Description)
	return task, err
}
