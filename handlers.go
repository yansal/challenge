package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/gin-gonic/gin.v1"
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
	authorized.PATCH("/tasks/:id", patchTasksIDHandler)
	authorized.POST("/tasks/:id/comments", postTasksIDCommentsHandler)
}

func getTasksHandler(c *gin.Context) {
	var tasks []TaskResource
	err := db.Select(&tasks, `SELECT tasks.id, tasks.created_at, tasks.updated_at, tasks.name, tasks.description, tasks.progression, users.id AS "user.id", users.username AS "user.username" FROM tasks JOIN users ON tasks.user_id = users.id ORDER BY created_at DESC;`)
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
	var task TaskResource
	err = db.Get(&task, `SELECT tasks.id, tasks.created_at, tasks.updated_at, tasks.name, tasks.description, tasks.progression, users.id AS "user.id", users.username AS "user.username" FROM tasks JOIN users ON tasks.user_id = users.id WHERE tasks.id = $1;`, id)
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
		log.Print("No user in POST /tasks/ handler context")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var task = Task{UserID: user.(User).ID}

	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if _, err := db.Exec(insertTask, task.Name, task.UserID, task.Description, task.Progression); err != nil {
		log.Printf("couldn't insert to tasks: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusCreated, "", nil)
}

type Patches []struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func patchTasksIDHandler(c *gin.Context) {
	user, exists := c.Get(gin.AuthUserKey)
	if !exists {
		log.Print("No user in PATCH /tasks/:id handler context")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var task TaskResource
	err = db.Get(&task, `SELECT tasks.id, tasks.created_at, tasks.updated_at, tasks.name, tasks.description, tasks.progression, users.id AS "user.id", users.username AS "user.username" FROM tasks JOIN users ON tasks.user_id = users.id WHERE tasks.id = $1;`, taskID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	if err != nil {
		log.Printf("couldn't select from tasks: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if task.User.ID != user.(User).ID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// TODO: Validate patch document
	var patches Patches
	if err := c.BindJSON(&patches); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// TODO: Apply patch document
	tx, err := db.Begin()
	if err != nil {
		log.Printf("couldn't begin: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	for _, patch := range patches {
		if patch.Op != "replace" {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf(`op must be "replace"; got %v`, patch.Op))
			return
		}
		if patch.Path != "/name" {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf(`path must be "/name"; got %v`, patch.Op))
			return
		}
		switch t := patch.Value.(type) {
		case string:
		default:
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf(`value must be a string; got %T`, t))
			return
		}
		if _, err := tx.Exec(`UPDATE tasks SET name = $1 WHERE id = $2;`, patch.Value, taskID); err != nil {
			log.Printf("couldn't update tasks: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		log.Printf("couldn't commit: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusNoContent, "", nil)
}

func getUsersIDTasksHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var tasks []TaskResource
	err = db.Select(&tasks, `SELECT tasks.id, tasks.created_at, tasks.updated_at, tasks.name, tasks.description, tasks.progression, users.id AS "user.id", users.username AS "user.username" FROM tasks JOIN users ON tasks.user_id = users.id WHERE user_id=$1 ORDER BY created_at DESC;`, id)
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
		log.Print("No user in POST /tasks/:id/comments handler context")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var i int
	err = db.Get(&i, "SELECT 1 FROM tasks WHERE id=$1", taskID)
	if err == sql.ErrNoRows {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("couldn't select from tasks: %v:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var comment = Comment{UserID: user.(User).ID, TaskID: taskID}
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if _, err := db.Exec(insertComment, comment.UserID, comment.TaskID, comment.Content); err != nil {
		log.Printf("couldn't insert to comments: %v", err)
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
	var comments []CommentResource
	err = db.Select(&comments, `SELECT comments.id, comments.created_at, comments.content, comments.task_id, users.id AS "user.id", users.username AS "user.username" FROM comments JOIN users ON comments.user_id = users.id WHERE task_id = $1 ORDER BY created_at DESC;`, taskID)
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
	var comments []CommentResource
	err = db.Select(&comments, `SELECT comments.id, comments.created_at, comments.content, comments.task_id, users.id AS "user.id", users.username AS "user.username" FROM comments JOIN users ON comments.user_id = users.id WHERE user_id = $1 ORDER BY created_at DESC;`, userID)
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
	var user User
	err := db.Get(&user, `SELECT * FROM users WHERE token=$1`, fields[1])
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
