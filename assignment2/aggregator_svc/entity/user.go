package entity

import (
	"time"
)

type UserResponse struct {
	ID        int32      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Balance   float64    `json:"balance"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type UserCreateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
