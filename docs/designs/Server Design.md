# Taskforge Server Design

## Scope

### Goals

#### MVP

- HTTP API Server that supports all the functions of a backend
- A Backend implementation for clients of this server
- The API server will use a taskforge.d config yaml and will then simply save
  to whatever the configured backend is.
- Token based authentication with tokens generated server side via a CLI command.
  The tokens will need to be stored in a backend independent way, a local state
  file seems like a good choice initially.

## Design

### Architecture

The basic flow for a "normal" client will look like this:

```text
Taskforge Client -> Taskforge Server Backend -> Taskforge Server API Server -> Taskforge Server API Server's Backend
```

In a sense Taskforge Server is just another client for the Taskforge Backends
however, it is special in that it is not meant for consumption by humans but by
other Taskforge clients.

### Authorization / Authentication

At this time there is no intent to allow multiple "users" on a Taskforge Server
it is expected that if you're using a Taskforge Server you want a shared task
list on some remote which can be reached by multiple computers.

That being said, a Taskforge Server will support multiple API Tokens. This will 
allow users to have a token per machine and if a machine or token becomes 
compromised they can simply revoke that token and not have to create a whole new
set of tokens.

Tokens will be generated with the command:

```bash
task server gen-token
```

This will generate an API Token and print it for the user. This then will be
stored in a plain text file `state.json`. This file will be in `~/.taskforge.d`
by default but will be configurable using the Taskforge `config.yaml`.

In future versions we will optionally encrypt, or support some way of doing
so, this file.

Clients will then authenticate using the HTTP `Authorization` header as follows:

```text
Authorization: Token $API_TOKEN
```

### API Endpoints

The server will have the following API endpoints with these purposes:

- `GET /list(?q=.*)`
  - Get all tasks in the list, if the q form parameter is given it will
    interpreted using the [Taskforge Query Language](https://github.com/chasinglogic/taskforge/blob/master/docs/design/Initial%20Design.md)
    returning the matching documents as a JSON array.
- `GET /list/current`
  - Return the current task.
- `POST /list`
  - Create multiple tasks, the request body should be a JSON array of tasks.
- `POST /task`
  - Create a single task
- `GET /task/$TASK_ID`
  - Get a single task by Task ID
- `PUT /task/$TASK_ID`
  - Update the task indicated by $TASK_ID, the request body should be a JSON Task
    and the response will be the updated Task as JSON.
- `PUT /task/$TASK_ID/complete`
  - Complete a task by ID. The response will be empty.
- `PUT /task/$TASK_ID/addNote`
  - Add a Note to the task indicated by ID. The request body should be a JSON
    Note. The response will be empty.

### Future Work

- Support websockets for notifications / reminders