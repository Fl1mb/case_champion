package models

import (
	"time"
)

type TaskModel struct {
	TaskID       int32     `json:"task_id"`
	FolderID     int32     `json:"folder_id"`
	UserID       int32     `json:"user_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Due_time     time.Time `json:"due_time"`
	Priority     int32     `json:"priority"`
	Is_completed bool      `json:"is_completed"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
}

type FolderModel struct {
	FolderID   int32     `json:"folder_id"`
	UserID     int32     `json:"user_id"`
	Name       string    `json:"name"`
	Created_at time.Time `json:"created_at"`
}
