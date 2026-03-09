package dto

import "dispatch/internal/platform/db"

type CreateUserRequest struct {
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password" binding:"required,min=8"`
}

type ListUsersParams struct {
	Pagination db.Pagination `json:"pagination"`
}
