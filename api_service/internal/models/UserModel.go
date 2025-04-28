package models

import "time"

type UserModel struct {
	ID         int32     `json:"user_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Fullname   string    `json:"full_name"`
	Is_active  bool      `json:"is_active"`
	Created_at time.Time `json:"created_at"`
}
