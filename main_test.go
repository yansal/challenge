package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestMain(m *testing.M) {
	flag.Parse()
	setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

var ts *httptest.Server

func setup() {
	ts = httptest.NewServer(router)

	db = sqlx.MustConnect("postgres", "dbname=challengetest sslmode=disable")
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
		db.MustExec(insertTask, task.Name, task.UserID, task.Description)
	}
	for _, comment := range []Comment{
		{TaskID: 1, UserID: 1, Content: "This is the first comment"},
		{TaskID: 1, UserID: 2, Content: "This is the second comment"},
		{TaskID: 2, UserID: 2, Content: "This is the third comment"},
	} {
		db.MustExec(insertComment, comment.UserID, comment.TaskID, comment.Content)
	}
}

func teardown() {
	db.MustExec(drop)
}

func TestGetTasks(t *testing.T) {
	resp, _ := http.Get(ts.URL + "/tasks/")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v; got %v", http.StatusOK, resp.StatusCode)
	}

	var tasks []TaskResource
	json.NewDecoder(resp.Body).Decode(&tasks)

	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks; got %d (%+v)", len(tasks), tasks)
	}
	if tasks[0].User.Username != "Alice" {
		t.Errorf("expected first task to embed username %q; got %q", "Alice", tasks[0].User.Username)
	}
}

func TestGetTasksID(t *testing.T) {
	resp, _ := http.Get(ts.URL + "/tasks/1")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v; got %v", http.StatusOK, resp.StatusCode)
	}

	var task TaskResource
	json.NewDecoder(resp.Body).Decode(&task)
	if task.ID != 1 {
		t.Errorf("expected id 1; got %d (%+v)", task.ID, task)
	}
	if task.User.Username != "Alice" {
		t.Errorf("expected username %q; got %q", "Alice", task.User.Username)
	}
}

func TestGetTasksIDNotFound(t *testing.T) {
	resp, _ := http.Get(ts.URL + "/tasks/123456")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code %v; got %v", http.StatusNotFound, resp.StatusCode)
	}
}

func TestGetTasksIDBadID(t *testing.T) {
	resp, _ := http.Get(ts.URL + "/tasks/hello")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %v; got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestPostTasks(t *testing.T) {
	name := "Posted task"
	marshalledTask, _ := json.Marshal(Task{Name: name})
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/tasks/", bytes.NewReader(marshalledTask))
	req.Header.Add("Authorization", "Token 077000ac559e1ba0fe4f303b614f30da6306341f")
	req.Header.Add("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status code %v; got %v", http.StatusCreated, resp.StatusCode)
	}

	resp, _ = http.Get(ts.URL + "/tasks/4")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v; got %v", http.StatusOK, resp.StatusCode)
	}

	var task TaskResource
	json.NewDecoder(resp.Body).Decode(&task)

	if task.User.Username != "Alice" {
		t.Errorf("expected username %q; got %q", "Alice", task.User.Username)
	}
	if task.Name != name {
		t.Errorf("expected name %q; got %q", name, task.Name)
	}
}

func TestPostTasksNoName(t *testing.T) {
	task, _ := json.Marshal(Task{})
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/tasks/", bytes.NewReader(task))
	req.Header.Add("Authorization", "Token 077000ac559e1ba0fe4f303b614f30da6306341f")
	req.Header.Add("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %v; got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestPostTasksUnauthenticated(t *testing.T) {
	marshalledTask, _ := json.Marshal(Task{Name: "Posted task"})
	resp, _ := http.Post(ts.URL+"/tasks/", "application/json", bytes.NewReader(marshalledTask))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %v; got %v", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestPostTasksBadAuthentication(t *testing.T) {
	marshalledTask, _ := json.Marshal(Task{Name: "Posted task"})
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/tasks/", bytes.NewReader(marshalledTask))
	req.Header.Add("Authorization", "Token 123456")
	req.Header.Add("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %v; got %v", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestGetUsersIDTasks(t *testing.T) {
	resp, _ := http.Get(ts.URL + "/users/1/tasks")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v; got %v", http.StatusOK, resp.StatusCode)
	}

	var tasks []TaskResource
	json.NewDecoder(resp.Body).Decode(&tasks)

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
	resp, _ := http.Get(ts.URL + "/users/123456/tasks")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v; got %v", http.StatusOK, resp.StatusCode)
	}

	var tasks []TaskResource
	json.NewDecoder(resp.Body).Decode(&tasks)
	if len(tasks) != 0 {
		t.Errorf("expected 0 task; got %d (%+v)", len(tasks), tasks)
	}
}

func TestGetUsersIDTasksBadID(t *testing.T) {
	resp, _ := http.Get(ts.URL + "/users/hello/tasks")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %v; got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestPostTasksIDComments(t *testing.T) {
	marshalledComment, _ := json.Marshal(Comment{Content: "Posted comment"})
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/tasks/1/comments", bytes.NewReader(marshalledComment))
	req.Header.Add("Authorization", "Token 077000ac559e1ba0fe4f303b614f30da6306341f")
	req.Header.Add("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status code %v; got %v", http.StatusCreated, resp.StatusCode)
	}
}

func TestPostTasksIDCommentsNoContent(t *testing.T) {
	marshalledComment, _ := json.Marshal(Comment{})
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/tasks/1/comments", bytes.NewReader(marshalledComment))
	req.Header.Add("Authorization", "Token 077000ac559e1ba0fe4f303b614f30da6306341f")
	req.Header.Add("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %v; got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestPostTasksIDCommentsDoesntExist(t *testing.T) {
	marshalledComment, _ := json.Marshal(Comment{Content: "This task doesn't exist"})
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/tasks/123456/comments", bytes.NewReader(marshalledComment))
	req.Header.Add("Authorization", "Token 077000ac559e1ba0fe4f303b614f30da6306341f")
	req.Header.Add("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code %v; got %v", http.StatusNotFound, resp.StatusCode)
	}
}

func TestPostTasksIDCommentsBadID(t *testing.T) {
	marshalledComment, _ := json.Marshal(Comment{Content: "This is not a task ID"})
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/tasks/hello/comments", bytes.NewReader(marshalledComment))
	req.Header.Add("Authorization", "Token 077000ac559e1ba0fe4f303b614f30da6306341f")
	req.Header.Add("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %v; got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestGetTasksIDComments(t *testing.T) {
	resp, _ := http.Get(ts.URL + "/tasks/1/comments")
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
	resp, _ := http.Get(ts.URL + "/users/1/comments")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v; got %v", http.StatusOK, resp.StatusCode)
	}

	var comments []CommentResource
	json.NewDecoder(resp.Body).Decode(&comments)

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

func TestPatchTasksIDForbidden(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/tasks/3", nil)
	req.Header.Add("Authorization", "Token 077000ac559e1ba0fe4f303b614f30da6306341f")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status code %v; got %v", http.StatusForbidden, resp.StatusCode)
	}
}

func TestPatchTasksIDDoesntExist(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/tasks/123456", nil)
	req.Header.Add("Authorization", "Token 077000ac559e1ba0fe4f303b614f30da6306341f")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code %v; got %v", http.StatusNotFound, resp.StatusCode)
	}
}

func TestPatchTasksBadID(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/tasks/hello", nil)
	req.Header.Add("Authorization", "Token 077000ac559e1ba0fe4f303b614f30da6306341f")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %v; got %v", http.StatusBadRequest, resp.StatusCode)
	}
}
