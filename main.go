package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"

	_ "github.com/lib/pq"
)

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
