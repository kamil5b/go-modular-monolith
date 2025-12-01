package noop

import (
	"go-modular-monolith/internal/domain/auth"
	"net/http"
)

type NoopHandler struct{}

func NewNoopHandler() *NoopHandler {
	return &NoopHandler{}
}

func (h *NoopHandler) Login(c auth.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "auth not implemented"})
}

func (h *NoopHandler) Register(c auth.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "auth not implemented"})
}

func (h *NoopHandler) Logout(c auth.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "auth not implemented"})
}

func (h *NoopHandler) RefreshToken(c auth.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "auth not implemented"})
}

func (h *NoopHandler) ValidateToken(c auth.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "auth not implemented"})
}

func (h *NoopHandler) ChangePassword(c auth.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "auth not implemented"})
}

func (h *NoopHandler) GetProfile(c auth.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "auth not implemented"})
}

func (h *NoopHandler) GetSessions(c auth.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "auth not implemented"})
}

func (h *NoopHandler) RevokeSession(c auth.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "auth not implemented"})
}

func (h *NoopHandler) RevokeAllSessions(c auth.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "auth not implemented"})
}
