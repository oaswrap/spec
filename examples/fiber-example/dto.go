package main

type LoginRequest struct {
	Username   string `json:"username" example:"john_doe" validate:"required"`
	Password   string `json:"password" example:"password123" validate:"required"`
	RememberMe bool   `json:"remember_me" example:"true"`
}

type RegisterRequest struct {
	Name                 string `json:"name" example:"John Doe" validate:"required"`
	Username             string `json:"username" example:"john_doe" validate:"required"`
	Email                string `json:"email" example:"john_doe@example.com" validate:"required,email"`
	Password             string `json:"password" example:"password123" validate:"required"`
	PasswordConfirmation string `json:"password_confirmation" example:"password123" validate:"required"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." validate:"required"`
}

type UpdateProfileRequest struct {
	Name     string `json:"name" example:"John Doe" validate:"required"`
	Email    string `json:"email" example:"john_doe@example.com" validate:"required,email"`
	Username string `json:"username" example:"john_doe" validate:"required"`
}

type UpdatePasswordRequest struct {
	CurrentPassword         string `json:"current_password" example:"current_password123" validate:"required"`
	NewPassword             string `json:"new_password" example:"new_password123" validate:"required"`
	NewPasswordConfirmation string `json:"new_password_confirmation" example:"new_password123" validate:"required"`
}

type DeleteAccountRequest struct {
	Password string `json:"password" example:"password123" validate:"required"`
}

type GetUserDetailRequest struct {
	UserID string `params:"user_id" path:"user_id" example:"12345"`
}

type Token struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type User struct {
	ID       string `json:"id" example:"12345"`
	Username string `json:"username" example:"john_doe"`
	Email    string `json:"email" example:"john_doe@example.com"`
	Name     string `json:"name" example:"John Doe"`
}

type Response[T any] struct {
	Status int `json:"status" example:"200"`
	Data   T   `json:"data,omitempty"`
}

type MessageResponse struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Operation successful"`
}

type ErrorResponse struct {
	Status int    `json:"status" example:"400"`
	Title  string `json:"title" example:"Bad Request"`
	Detail string `json:"detail" example:"Invalid input data"`
}

type ValidationResponse struct {
	Status int          `json:"status" example:"422"`
	Title  string       `json:"title" example:"Validation Error"`
	Detail string       `json:"detail" example:"Input data does not meet validation criteria"`
	Errors []FieldError `json:"errors" example:"[{\"field\":\"email\",\"message\":\"Email is required\"}]"`
}

type FieldError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Email is required"`
}
