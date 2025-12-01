package v1

import (
	"go-modular-monolith/internal/domain/auth"
	"net/http"
)

type Handler struct {
	svc auth.AuthService
}

func NewHandler(s auth.AuthService) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) Login(c auth.Context) error {
	var req auth.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userAgent := c.GetUserAgent()
	ipAddress := c.GetClientIP()

	resp, err := h.svc.Login(c.GetContext(), &req, userAgent, ipAddress)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) Register(c auth.Context) error {
	var req auth.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.svc.Register(c.GetContext(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) Logout(c auth.Context) error {
	userID := c.GetUserID()
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req auth.LogoutRequest
	_ = c.Bind(&req)

	if err := h.svc.Logout(c.GetContext(), userID, &req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, auth.MessageResponse{Message: "Logged out successfully", Success: true})
}

func (h *Handler) RefreshToken(c auth.Context) error {
	var req auth.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.svc.RefreshToken(c.GetContext(), req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) ValidateToken(c auth.Context) error {
	var req auth.ValidateTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.svc.ValidateToken(c.GetContext(), req.Token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) ChangePassword(c auth.Context) error {
	userID := c.GetUserID()
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req auth.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.svc.ChangePassword(c.GetContext(), userID, &req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, auth.MessageResponse{Message: "Password changed successfully", Success: true})
}

func (h *Handler) GetProfile(c auth.Context) error {
	userID := c.GetUserID()
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	resp, err := h.svc.ValidateToken(c.GetContext(), token)
	if err != nil || !resp.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	return c.JSON(http.StatusOK, resp.User)
}

func (h *Handler) GetSessions(c auth.Context) error {
	userID := c.GetUserID()
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	resp, err := h.svc.GetSessions(c.GetContext(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) RevokeSession(c auth.Context) error {
	userID := c.GetUserID()
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	sessionID := c.Param("id")
	if sessionID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "session id required"})
	}

	if err := h.svc.RevokeSession(c.GetContext(), userID, sessionID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, auth.MessageResponse{Message: "Session revoked successfully", Success: true})
}

func (h *Handler) RevokeAllSessions(c auth.Context) error {
	userID := c.GetUserID()
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	if err := h.svc.RevokeAllSessions(c.GetContext(), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, auth.MessageResponse{Message: "All sessions revoked successfully", Success: true})
}
