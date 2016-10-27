# Task Manager

[![Build Status](https://travis-ci.org/yansal/task-manager.svg?branch=master)](https://travis-ci.org/yansal/task-manager)

This is a task manager API with the following endpoints:

* GET /tasks/
* GET /tasks/:id
* POST /tasks/
* PATCH /tasks/:id
* GET /users/:id/tasks
* POST /tasks/:id/comments
* GET /tasks/:id/comments
* GET /users/:id/comments

A public instance is reachable at https://yansal-task-manager.herokuapp.com.

## Run

Get the code with `go get github.com/yansal/task-manager`.

A PostgreSQL instance is required. On macOS, install it with `brew install postgresql` and start it with `brew services start postgresql`

When running for the first time, it is necessary to seed the database. To do so, run `$GOPATH/bin/task-manager -seed`. Otherwise, just run `$GOPATH/bin/task-manager`. The default database is `$USERNAME` and the default listening port is 8080. Both can be configured with env (respectively `DATABASE_URL` and `PORT`).

Run the test suite with `go test github.com/yansal/task-manager`. The testing database name is `taskmanagertest` and must be created with `createdb taskmanagertest`.
