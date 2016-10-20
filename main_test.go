package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if err := setup(); err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

var ts *httptest.Server

func setup() error {
	ts = httptest.NewServer(router)

	var err error
	if db, err = sqlx.Connect("postgres", "dbname=challengetest sslmode=disable"); err != nil {
		return fmt.Errorf("couldn't connect to database: %v", err)
	}
	db.MustExec(create)
	if err := seedTableUsers(); err != nil {
		return err
	}
	if err := seedTableTasks(); err != nil {
		return err
	}
	if err := seedTableComments(); err != nil {
		return err
	}
	return nil
}

func seedTableUsers() error {
	for _, user := range []User{
		{Username: "Alice", Token: "077000ac559e1ba0fe4f303b614f30da6306341f"},
		{Username: "Bob", Token: "ef2e253a2b4564ae949b053025c845552f2e99cc"},
	} {
		if err := insertUser(user); err != nil {
			return err
		}
	}
	return nil
}

func seedTableTasks() error {
	for _, task := range []Task{
		{Name: "First task", UserID: 1, Description: "This is the first task"},
		{Name: "Second task", UserID: 1, Description: "This is the second task"},
		{Name: "Third task", UserID: 2, Description: "This is the third task"},
	} {
		if err := insertTask(task); err != nil {
			return err
		}
	}
	return nil
}

func seedTableComments() error {
	for _, comment := range []Comment{
		{TaskID: 1, UserID: 1, Content: "This is the first comment"},
		{TaskID: 1, UserID: 2, Content: "This is the second comment"},
		{TaskID: 2, UserID: 2, Content: "This is the third comment"},
	} {
		if err := insertComment(comment); err != nil {
			return err
		}
	}
	return nil
}

func teardown() {
	db.MustExec(`DROP TABLE comments; DROP TABLE tasks; DROP TABLE users;`)
}

func TestGetTasks(t *testing.T) {
	resp, err := http.Get(ts.URL + "/tasks/")
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v; got %v", http.StatusOK, resp.StatusCode)
	}
	var tasks []TaskResource
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		t.Errorf("couldn't decode body to JSON: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks; got %d (%+v)", len(tasks), tasks)
	}
	if tasks[0].User.Username != "Alice" {
		t.Errorf("expected first task to embed username %q; got %q", "Alice", tasks[0].User.Username)
	}
}

func TestGetTasksID(t *testing.T) {
	resp, err := http.Get(ts.URL + "/tasks/1")
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v; got %v", http.StatusOK, resp.StatusCode)
	}
	var task TaskResource
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		t.Errorf("couldn't decode body to JSON: %v", err)
	}
	if task.ID != 1 {
		t.Errorf("expected id 1; got %d (%+v)", task.ID, task)
	}
	if task.User.Username != "Alice" {
		t.Errorf("expected username %q; got %q", "Alice", task.User.Username)
	}
}

func TestGetTasksIDNotFound(t *testing.T) {
	resp, err := http.Get(ts.URL + "/tasks/123456")
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code %v; got %v", http.StatusNotFound, resp.StatusCode)
	}
}

func TestGetTasksIDBadID(t *testing.T) {
	resp, err := http.Get(ts.URL + "/tasks/hello")
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %v; got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestPostTasks(t *testing.T) {
	name := "Posted task"
	marshalledTask, err := json.Marshal(Task{Name: name})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", ts.URL+"/tasks/", bytes.NewReader(marshalledTask))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Token 077000ac559e1ba0fe4f303b614f30da6306341f")
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status code %v; got %v", http.StatusCreated, resp.StatusCode)
	}

	resp, err = http.Get(ts.URL + "/tasks/4")
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v; got %v", http.StatusOK, resp.StatusCode)
	}
	var task TaskResource
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		t.Errorf("couldn't decode body to JSON: %v", err)
	}
	if task.User.Username != "Alice" {
		t.Errorf("expected username %q; got %q", "Alice", task.User.Username)
	}
	if task.Name != name {
		t.Errorf("expected name %q; got %q", name, task.Name)
	}
}

func TestPostTasksNoName(t *testing.T) {
	task, err := json.Marshal(Task{})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", ts.URL+"/tasks/", bytes.NewReader(task))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Token 077000ac559e1ba0fe4f303b614f30da6306341f")
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %v; got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestPostTasksUnauthenticated(t *testing.T) {
	marshalledTask, err := json.Marshal(Task{Name: "Posted task"})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(ts.URL+"/tasks/", "application/json", bytes.NewReader(marshalledTask))
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %v; got %v", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestPostTasksBadAuthentication(t *testing.T) {
	marshalledTask, err := json.Marshal(Task{Name: "Posted task"})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", ts.URL+"/tasks/", bytes.NewReader(marshalledTask))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Token 123456")
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %v; got %v", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestGetUsersIDTasks(t *testing.T) {
	resp, err := http.Get(ts.URL + "/users/1/tasks")
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v; got %v", http.StatusOK, resp.StatusCode)
	}

	var tasks []TaskResource
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		t.Errorf("couldn't decode body to JSON: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks; got %d (%+v)", len(tasks), tasks)
	}
	for _, task := range tasks {
		if task.User.Username != "Alice" {
			t.Errorf("expected username %q; got %q", "Alice", task.User.Username)
		}
	}
}

func TestGetUsersIDTasksDoesntExist(t *testing.T) {
	resp, err := http.Get(ts.URL + "/users/123456/tasks")
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v; got %v", http.StatusOK, resp.StatusCode)
	}
	var tasks []TaskResource
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		t.Errorf("couldn't decode body to JSON: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 task; got %d (%+v)", len(tasks), tasks)
	}
}

func TestGetUsersIDTasksBadID(t *testing.T) {
	resp, err := http.Get(ts.URL + "/users/hello/tasks")
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %v; got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestPostTasksIDComments(t *testing.T) {
	marshalledComment, err := json.Marshal(Comment{Content: "Posted comment"})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", ts.URL+"/tasks/1/comments", bytes.NewReader(marshalledComment))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Token 077000ac559e1ba0fe4f303b614f30da6306341f")
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status code %v; got %v", http.StatusCreated, resp.StatusCode)
	}
}

func TestPostTasksIDCommentsNoContent(t *testing.T) {
	marshalledComment, err := json.Marshal(Comment{})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", ts.URL+"/tasks/1/comments", bytes.NewReader(marshalledComment))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Token 077000ac559e1ba0fe4f303b614f30da6306341f")
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %v; got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestPostTasksIDCommentsDoesntExist(t *testing.T) {
	marshalledComment, err := json.Marshal(Comment{Content: "This task doesn't exist"})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", ts.URL+"/tasks/123456/comments", bytes.NewReader(marshalledComment))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Token 077000ac559e1ba0fe4f303b614f30da6306341f")
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		t.Errorf("unexpected status code %v", http.StatusInternalServerError)
	}
}

func TestGetTasksIDComments(t *testing.T) {
	resp, err := http.Get(ts.URL + "/tasks/1/comments")
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v; got %v", http.StatusOK, resp.StatusCode)
	}

	var comments []CommentResource
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		t.Errorf("couldn't decode body to JSON: %v", err)
	}
	if len(comments) != 3 {
		t.Errorf("expected 3 comments; got %d (%+v)", len(comments), comments)
	}
	for _, comment := range comments {
		if comment.TaskID != 1 {
			t.Errorf("expected task_id 1; got %d", comment.TaskID)
		}
	}
	if len(comments) == 0 {
		t.FailNow()
	}
	if comments[0].User.Username != "Alice" {
		t.Errorf("expected first comment to embed username %q; got %q", "Alice", comments[0].User.Username)
	}
}

func TestGetUsersIDComments(t *testing.T) {
	resp, err := http.Get(ts.URL + "/users/1/comments")
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v; got %v", http.StatusOK, resp.StatusCode)
	}

	var comments []CommentResource
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		t.Errorf("couldn't decode body to JSON: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("expected 2 comments; got %d (%+v)", len(comments), comments)
	}
	for _, comment := range comments {
		if comment.User.Username != "Alice" {
			t.Errorf("expected username %q; got %q", "Alice", comment.User.Username)
		}
	}
	if len(comments) == 0 {
		t.FailNow()
	}
	if comments[0].TaskID != 1 {
		t.Errorf("expected first comment to have task_id 1; got %d", comments[0].TaskID)
	}
}
