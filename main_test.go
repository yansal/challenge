package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if err := setup(); err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	if err := teardown(); err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

var ts *httptest.Server

func setup() error {
	ts = httptest.NewServer(router)

	var err error
	if db, err = sql.Open("postgres", "dbname=challengetest sslmode=disable"); err != nil {
		return fmt.Errorf("couldn't open database connection: %v", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("couldn't ping database: %v", err)
	}
	if err := createTableUsers(); err != nil {
		return err
	}
	if err := createTableTasks(); err != nil {
		return err
	}
	if err := seedTableUsers(); err != nil {
		return err
	}
	if err := seedTableTasks(); err != nil {
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

func teardown() error {
	if _, err := db.Exec("DROP TABLE tasks;"); err != nil {
		return fmt.Errorf("couldn't drop tasks table: %v", err)
	}
	if _, err := db.Exec("DROP TABLE users;"); err != nil {
		return fmt.Errorf("couldn't drop tasks table: %v", err)
	}
	return nil
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

	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		t.Errorf("couldn't decode body to JSON: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks; got %d", len(tasks))
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
	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		t.Errorf("couldn't decode body to JSON: %v", err)
	}
	if task.ID != 1 {
		t.Errorf("expected id 1; got %v", task.ID)
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

func TestGetTasksIDBadRequest(t *testing.T) {
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
	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		t.Errorf("couldn't decode body to JSON: %v", err)
	}
	if task.UserID != 1 {
		t.Errorf("expected user_id 1; got %v", task.UserID)
	}
	if task.Name != name {
		t.Errorf("expected name %q; got %q", name, task.Name)
	}
}

func TestPostTasksBadRequest(t *testing.T) {
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

func TestPatchTasksID(t *testing.T) {
	marshalledTask, err := json.Marshal(Task{Name: "New name for this task"})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PATCH", ts.URL+"/tasks/1", bytes.NewReader(marshalledTask))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status code %v; got %v", http.StatusNoContent, resp.StatusCode)
	}
}
