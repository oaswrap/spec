package dto

import "time"

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ErrorResponse struct {
	Status int    `json:"status" example:"400"`
	Title  string `json:"title" example:"Bad Request"`
	Detail string `json:"detail,omitempty" example:"Invalid input provided"`
}

type ValidationResponse struct {
	Status int          `json:"status" example:"422"`
	Title  string       `json:"title" example:"Validation Error"`
	Detail string       `json:"detail,omitempty" example:"Input validation failed"`
	Errors []FieldError `json:"errors,omitempty"`
}

type FieldError struct {
	Field   string `json:"field" example:"username"`
	Message string `json:"message" example:"Username is required"`
}

type Response[T any] struct {
	Status int `json:"status" example:"200"`
	Data   T   `json:"data"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProfile struct {
	ID              string     `json:"id"`
	Username        string     `json:"username"`
	Email           string     `json:"email"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	FullName        string     `json:"full_name"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
