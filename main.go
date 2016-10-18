package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/gin-gonic/gin.v1"

	_ "github.com/lib/pq"
)

func getTasksHandler(c *gin.Context) {
	tasks, err := selectTasks()
	if err != nil {
		log.Printf("couldn't select from tasks: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
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
		log.Printf("couldn't select from tasks: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, task)
}

func postTasksHandler(c *gin.Context) {
	user, exists := c.Get(gin.AuthUserKey)
	if !exists {
		log.Print("No user in POST /tasks/ handler")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var task = Task{UserID: user.(User).ID}

	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if err := insertTask(task); err != nil {
		log.Printf("couldn't insert into tasks: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusCreated, "", nil)
}

func getUsersIDTasksHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	tasks, err := selectTasksWhereUser(id)
	if err != nil {
		log.Printf("couldn't select from tasks: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func postTasksIDCommentsHandler(c *gin.Context) {
	user, exists := c.Get(gin.AuthUserKey)
	if !exists {
		log.Print("No user in POST /tasks/:id/comments handler")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var comment = Comment{UserID: user.(User).ID, TaskID: taskID}
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if err := insertComment(comment); err != nil {
		log.Printf("couldn't insert into comments: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusCreated, "", nil)
}

func getTasksIDCommentsHandler(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	comments, err := selectCommentsWhereTask(taskID)
	if err != nil {
		log.Printf("couldn't select from comments: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, comments)
}

func getUsersIDCommentsHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	comments, err := selectCommentsWhereUser(userID)
	if err != nil {
		log.Printf("couldn't select from tasks: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, comments)
}

func authMiddleware(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	fields := strings.Fields(header)
	if len(fields) != 2 || fields[0] != "Token" {
		c.Header("WWW-Authenticate", "Token")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user, err := selectUser(fields[1])
	if err == sql.ErrNoRows {
		c.Header("WWW-Authenticate", "Token")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Printf("couldn't select from users: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Set(gin.AuthUserKey, user)
}

var (
	db     *sql.DB
	router *gin.Engine
)

func init() {
	router = gin.Default()
	router.GET("/tasks/", getTasksHandler)
	router.GET("/tasks/:id", getTasksIDHandler)
	router.GET("/users/:id/tasks", getUsersIDTasksHandler)
	router.GET("/tasks/:id/comments", getTasksIDCommentsHandler)
	router.GET("/users/:id/comments", getUsersIDCommentsHandler)
	authorized := router.Group("/", authMiddleware)
	authorized.POST("/tasks/", postTasksHandler)
	authorized.POST("/tasks/:id/comments", postTasksIDCommentsHandler)
}

func main() {
	var err error
	db, err = sql.Open("postgres", "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(router.Run())
}
