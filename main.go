package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/gin-gonic/gin.v1"

	_ "github.com/lib/pq"
)

type Task struct {
	ID           int        `json:"id,omitempty"`
	Name         string     `json:"name,omitempty" binding:"required"`
	User         string     `json:"user,omitempty" binding:"required"`
	Description  string     `json:"description,omitempty"`
	CreationDate *time.Time `json:"creation_date,omitempty"`
}

func createTableTasks() error {
	if _, err := db.Exec(`CREATE TABLE tasks(
    "id" SERIAL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "user" TEXT NOT NULL,
    "description" TEXT,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`); err != nil {
		return fmt.Errorf(`couldn't create "tasks" table: %v`, err)
	}
	return nil
}

func insertTask(task Task) error {
	if _, err := db.Exec(`INSERT INTO tasks ("name", "user", "description") VALUES ($1, $2, $3)`, task.Name, task.User, task.Description); err != nil {
		return fmt.Errorf("couldn't insert into tasks: %v", err)
	}
	return nil
}

func selectTasks() ([]Task, error) {
	rows, err := db.Query(`SELECT * FROM tasks ORDER BY "creation_date";`)
	if err != nil {
		return nil, fmt.Errorf("couldn't select from tasks table: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Name, &task.User, &task.Description, &task.CreationDate); err != nil {
			return tasks, fmt.Errorf(`couldn't scan row: %v`, err)
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func selectTask(id int) (Task, error) {
	row := db.QueryRow(`SELECT * FROM tasks WHERE id=$1`, id)
	var task Task
	err := row.Scan(&task.ID, &task.Name, &task.User, &task.Description, &task.CreationDate)
	return task, err
}

func getTasksHandler(c *gin.Context) {
	tasks, err := selectTasks()
	if err != nil {
		log.Print(err)
		c.Data(http.StatusInternalServerError, "", nil)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func getTasksIDHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	task, err := selectTask(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	if err != nil {
		log.Print(err)
		c.Data(http.StatusInternalServerError, "", nil)
		return
	}
	c.JSON(http.StatusOK, task)
}

func postTasksHandler(c *gin.Context) {
	var task Task
	if err := c.Bind(&task); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if err := insertTask(task); err != nil {
		log.Print(err)
		c.Data(http.StatusInternalServerError, "", nil)
		return
	}
	c.Data(http.StatusCreated, "", nil)
}

var (
	db     *sql.DB
	router *gin.Engine
)

func init() {
	router = gin.Default()
	router.GET("/tasks/", getTasksHandler)
	router.GET("/tasks/:id", getTasksIDHandler)
	router.POST("/tasks/", postTasksHandler)
}

func main() {
	var err error
	db, err = sql.Open("postgres", "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(router.Run())
}
