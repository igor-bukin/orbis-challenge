package models

import (
	"time"

	"github.com/google/uuid"
)

// User model.
type User struct {
	tableName struct{} `sql:"users" pg:",discard_unknown_columns"` // nolint

	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time
}

// LoginRequest login request model.
type LoginRequest struct {
	Email    string `json:"email"      validate:"required|email"`
	Password string `json:"password"   validate:"required"`
}

// SignUpRequest signup request model.
type SignUpRequest struct {
	Email    string `json:"login"      validate:"required|email"`
	Password string `json:"password"   validate:"required"`
}
