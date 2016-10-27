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

## Example of use

All examples use [HTTPie](https://httpie.org)

### Get all tasks
```
$ http https://yansal-task-manager.herokuapp.com/tasks/
HTTP/1.1 200 OK
Connection: keep-alive
Content-Length: 615
Content-Type: application/json; charset=utf-8
Date: Thu, 27 Oct 2016 22:31:06 GMT
Server: Cowboy
Via: 1.1 vegur

[
    {
        "created_at": "2016-10-27T22:30:56.185261Z", 
        "description": "This is the third task", 
        "id": 3, 
        "name": "Third task", 
        "progression": 0, 
        "updated_at": "2016-10-27T22:30:56.185261Z", 
        "user": {
            "id": 2, 
            "username": "Bob"
        }
    }, 
    {
        "created_at": "2016-10-27T22:30:56.183708Z", 
        "description": "This is the second task", 
        "id": 2, 
        "name": "Second task", 
        "progression": 0, 
        "updated_at": "2016-10-27T22:30:56.183708Z", 
        "user": {
            "id": 1, 
            "username": "Alice"
        }
    }, 
    {
        "created_at": "2016-10-27T22:30:56.18166Z", 
        "description": "This is the first task", 
        "id": 1, 
        "name": "First task", 
        "progression": 0, 
        "updated_at": "2016-10-27T22:30:56.18166Z", 
        "user": {
            "id": 1, 
            "username": "Alice"
        }
    }
]
```

### Get one task

This request returns an `Etag` header that is required for the PATCH /tasks/:id endpoint.
```
$ http https://yansal-task-manager.herokuapp.com/tasks/1
HTTP/1.1 200 OK
Connection: keep-alive
Content-Length: 615
Content-Type: application/json; charset=utf-8
Date: Thu, 27 Oct 2016 22:31:06 GMT
Server: Cowboy
Via: 1.1 vegur

{
    "created_at": "2016-10-27T22:30:56.18166Z", 
    "description": "This is the first task", 
    "id": 1, 
    "name": "First task", 
    "progression": 0, 
    "updated_at": "2016-10-27T22:30:56.18166Z", 
    "user": {
        "id": 1, 
        "username": "Alice"
    }
}
```
### Post a task
A `Authorization` header is required for every POST and PATCH endpoints.

To post a task, the `name` field is required and the `description` and `progression` fields are optional.
```
$ http -v https://yansal-task-manager.herokuapp.com/tasks/ name="Posted Task" Authorization:"Token 077000ac559e1ba0fe4f303b614f30da6306341f"
POST /tasks/ HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Authorization: Token 077000ac559e1ba0fe4f303b614f30da6306341f
Connection: keep-alive
Content-Length: 23
Content-Type: application/json
Host: yansal-task-manager.herokuapp.com
User-Agent: HTTPie/0.9.6

{
    "name": "Posted Task"
}

HTTP/1.1 201 Created
Connection: keep-alive
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Thu, 27 Oct 2016 22:36:21 GMT
Server: Cowboy
Via: 1.1 vegur
```

### Patch a task
To patch a task, there are two additional conditions:

1. The body must be a [JSON Patch document](http://jsonpatch.com/) (and the `Content-Type` must be `application/json-patch+json`).
2. A `If-Match` header must be present and match the task's `Etag`.
```
$ echo '[{"op":"replace","path":"/name","value":"Patched Task"}]' | http -v PATCH https://yansal-task-manager.herokuapp.com/tasks/1 Authorization:"Token 077000ac559e1ba0fe4f303b614f30da6306341f" Content-Type:application/json-patch+json If-Match:1ec62297510e447eef7729860773879a
PATCH /tasks/1 HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Authorization: Token 077000ac559e1ba0fe4f303b614f30da6306341f
Connection: keep-alive
Content-Length: 57
Content-Type: application/json-patch+json
Host: yansal-task-manager.herokuapp.com
If-Match: 1ec62297510e447eef7729860773879a
User-Agent: HTTPie/0.9.6

[
    {
        "op": "replace", 
        "path": "/name", 
        "value": "Patched Task"
    }
]

HTTP/1.1 204 No Content
Connection: keep-alive
Content-Length: 0
Date: Thu, 27 Oct 2016 22:39:07 GMT
Server: Cowboy
Via: 1.1 vegur
```
