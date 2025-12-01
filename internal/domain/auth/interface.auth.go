package auth

import (
	"context"
)

// Context defines the interface for HTTP context used by auth handlers
type Context interface {
	BindJSON(obj any) error
	BindURI(obj any) error
	BindQuery(obj any) error
	BindHeader(obj any) error
	Bind(obj any) error
	JSON(code int, v any) error
	Param(name string) string
	GetUserID() string
	Get(key string) any
	Set(key string, value any)
	GetContext() context.Context
	GetHeader(key string) string
	SetHeader(key, value string)
	GetCookie(name string) (string, error)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	RemoveCookie(name string)
	GetClientIP() string
	GetUserAgent() string
}

// AuthHandler defines the interface for authentication HTTP handlers
type AuthHandler interface {
	Login(c Context) error
	Register(c Context) error
	Logout(c Context) error
	RefreshToken(c Context) error
	ValidateToken(c Context) error
	ChangePassword(c Context) error
	GetProfile(c Context) error
	GetSessions(c Context) error
	RevokeSession(c Context) error
	RevokeAllSessions(c Context) error
}

// AuthService defines the interface for authentication business logic
type AuthService interface {
	// Authentication
	Login(ctx context.Context, req *LoginRequest, userAgent, ipAddress string) (*LoginResponse, error)
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Logout(ctx context.Context, userID string, req *LogoutRequest) error
	RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error)
	ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error)

	// Password management
	ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) error
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
	ConfirmResetPassword(ctx context.Context, req *ConfirmResetPasswordRequest) error

	// Session management
	GetSessions(ctx context.Context, userID string) (*SessionListResponse, error)
	RevokeSession(ctx context.Context, userID, sessionID string) error
	RevokeAllSessions(ctx context.Context, userID string) error

	// Token utilities
	GenerateAccessToken(claims *TokenClaims) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ParseToken(token string) (*TokenClaims, error)

	// Password utilities
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}

// AuthRepository defines the interface for authentication data access
type AuthRepository interface {
	StartContext(ctx context.Context) context.Context
	DeferErrorContext(ctx context.Context, err error)

	// Credential operations
	CreateCredential(ctx context.Context, cred *Credential) error
	GetCredentialByUsername(ctx context.Context, username string) (*Credential, error)
	GetCredentialByEmail(ctx context.Context, email string) (*Credential, error)
	GetCredentialByUserID(ctx context.Context, userID string) (*Credential, error)
	UpdateCredential(ctx context.Context, cred *Credential) error
	UpdatePassword(ctx context.Context, userID, passwordHash string) error
	UpdateLastLogin(ctx context.Context, userID string) error

	// Session operations
	CreateSession(ctx context.Context, session *Session) error
	GetSessionByToken(ctx context.Context, token string) (*Session, error)
	GetSessionByID(ctx context.Context, id string) (*Session, error)
	GetSessionsByUserID(ctx context.Context, userID string) ([]Session, error)
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeAllUserSessions(ctx context.Context, userID string) error
	DeleteExpiredSessions(ctx context.Context) error
}

// Middleware defines the interface for authentication middleware
type Middleware interface {
	// Authenticate validates the request and sets auth context
	Authenticate() func(next func(Context) error) func(Context) error

	// RequireAuth ensures the request is authenticated
	RequireAuth() func(next func(Context) error) func(Context) error

	// OptionalAuth tries to authenticate but allows unauthenticated requests
	OptionalAuth() func(next func(Context) error) func(Context) error

	// RequireRoles ensures the authenticated user has specific roles
	RequireRoles(roles ...string) func(next func(Context) error) func(Context) error
}
