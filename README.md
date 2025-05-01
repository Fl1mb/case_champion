# Task Management API Documentation

## Authentication
Protected routes require:
## Endpoints

### Authentication

#### Register User
`POST /register`
```json
Request:
{
  "username": "string",
  "email": "string",
  "password": "string",
  "fullname": "string"
}

Response (201):
{
  "user_id": "int32",
  "username": "string",
  "email": "string",
  "full_name": "string"
}

#### Login User

POST /login
json

Request:
{
  "username": "string",
  "password": "string"
}

Response:
{
  "access_token": "string",
  "expires_at": "string"
}

Folders
Get All Folders

GET /folders
json

Response:
[
  {
    "folder_id": "int32",
    "name": "string",
    "created_at": "string",
    "task_ids": ["int32"]
  }
]

Create Folder

POST /folders
json

Request:
{
  "name": "string"
}

Response (201):
{
  "folder_id": "int32",
  "name": "string",
  "created_at": "string"
}


Get Folder

GET /folders/{folderID}
json

Response:
{
  "folder_id": "int32",
  "name": "string",
  "created_at": "string",
  "task_ids": ["int32"]
}

Update Folder

PUT /folders/{folderID}
json

Request:
{
  "name": "string"
}

Response:
{
  "folder_id": "int32",
  "name": "string",
  "created_at": "string"
}

Delete Folder

DELETE /folders/{folderID}
json

Response:
{
  "success": "bool",
  "message": "string"
}


### Tasks
Get All Tasks

GET /tasks
json

Query Params:
?folder_id=int32 (optional)

Response:
{
  "tasks": [
    {
      "task_id": "int32",
      "title": "string",
      "description": "string",
      "due_time": "string",
      "priority": "int32",
      "is_completed": "bool"
    }
  ],
  "total_count": "int32"
}

Create Task

POST /tasks
json

Request:
{
  "folder_id": "int32",
  "title": "string",
  "description": "string",
  "due_time": "string",
  "priority": "int32"
}

Response (201):
{
  "task_id": "int32",
  "title": "string",
  "description": "string",
  "due_time": "string",
  "priority": "int32",
  "is_completed": "bool"
}

Get Task

GET /tasks/{taskID}
json

Response:
{
  "task_id": "int32",
  "title": "string",
  "description": "string",
  "due_time": "string",
  "priority": "int32",
  "is_completed": "bool"
}

Update Task

PUT /tasks/{taskID}
json

Request:
{
  "title": "string",
  "description": "string",
  "due_time": "string",
  "priority": "int32"
}

Response:
{
  "task_id": "int32",
  "title": "string",
  "description": "string",
  "due_time": "string",
  "priority": "int32",
  "is_completed": "bool"
}


Delete Task

DELETE /tasks/{taskID}
json

Response:
{
  "success": "bool",
  "message": "string"
}

Toggle Task Completion

PATCH /tasks/{taskID}/toggle
json

Response:
{
  "task_id": "int32",
  "is_completed": "bool"
}


Move Task

PATCH /tasks/{taskID}/move
json

Request:
{
  "new_folder_id": "int32"
}

Response:
{
  "task_id": "int32",
  "folder_id": "int32"
}

Search Tasks

GET /tasks/search
json

Query Params:
?query=string
?completed=bool
?priority=int32
?due_before=string

Response:
{
  "tasks": [
    {
      "task_id": "int32",
      "title": "string",
      "description": "string",
      "due_time": "string",
      "priority": "int32",
      "is_completed": "bool"
    }
  ],
  "total_count": "int32"
}

Error Responses
json

{
  "error": "Error message",
  "code": "HTTP status code"
}