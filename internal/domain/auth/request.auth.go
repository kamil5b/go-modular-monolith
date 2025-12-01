package auth

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// RefreshTokenRequest represents the refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest represents the change password request payload
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResetPasswordRequest represents the reset password request payload
type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ConfirmResetPasswordRequest represents the confirm reset password request payload
type ConfirmResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// LogoutRequest represents the logout request payload
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
	AllDevices   bool   `json:"all_devices,omitempty"`
}

// ValidateTokenRequest represents the validate token request payload
type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required"`
}
