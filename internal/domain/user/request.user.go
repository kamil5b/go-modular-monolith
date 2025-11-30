package user

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
	ID    string `json:"id" binding:"required"`
	Name  string `json:"name"`
	Email string `json:"email" binding:"omitempty,email"`
}
