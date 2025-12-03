package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"go-modular-monolith/internal/infrastructure/worker"
	userdomain "go-modular-monolith/internal/modules/user/domain"
	"go-modular-monolith/internal/shared/email"
)

// UserWorkerHandler processes user-related tasks
type UserWorkerHandler struct {
	userRepository userdomain.Repository
	emailService   email.EmailService
}

// NewUserWorkerHandler creates a new user worker handler
func NewUserWorkerHandler(userRepository userdomain.Repository, emailService email.EmailService) *UserWorkerHandler {
	return &UserWorkerHandler{
		userRepository: userRepository,
		emailService:   emailService,
	}
}

// HandleSendWelcomeEmail handles the welcome email task
func (h *UserWorkerHandler) HandleSendWelcomeEmail(ctx context.Context, payload worker.TaskPayload) error {
	var p SendWelcomeEmailPayload

	// Unmarshal payload
	data, _ := json.Marshal(payload)
	if err := json.Unmarshal(data, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Validate payload
	if p.UserID == "" || p.Email == "" {
		return fmt.Errorf("missing required fields in payload")
	}

	// Get user details
	user, err := h.userRepository.GetByID(ctx, p.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found: %s", p.UserID)
	}

	// Send welcome email
	emailMsg := &email.Email{
		To:       []string{user.Email},
		Subject:  "Welcome to Our Platform!",
		HTMLBody: fmt.Sprintf("<h1>Welcome %s!</h1><p>Thank you for joining us. We're excited to have you on board.</p>", user.Name),
		TextBody: fmt.Sprintf("Welcome %s!\n\nThank you for joining us. We're excited to have you on board.", user.Name),
	}

	if err := h.emailService.Send(ctx, emailMsg); err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	fmt.Printf("Sent welcome email to %s (%s)\n", user.Name, user.Email)

	return nil
}

// HandleSendPasswordResetEmail handles the password reset email task
func (h *UserWorkerHandler) HandleSendPasswordResetEmail(ctx context.Context, payload worker.TaskPayload) error {
	var p SendPasswordResetEmailPayload

	data, _ := json.Marshal(payload)
	if err := json.Unmarshal(data, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if p.UserID == "" || p.Email == "" || p.ResetLink == "" {
		return fmt.Errorf("missing required fields in payload")
	}

	user, err := h.userRepository.GetByID(ctx, p.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found: %s", p.UserID)
	}

	// Send password reset email
	emailMsg := &email.Email{
		To:       []string{user.Email},
		Subject:  "Password Reset Request",
		HTMLBody: fmt.Sprintf("<h1>Password Reset</h1><p>Click <a href=\"%s\">here</a> to reset your password.</p>", p.ResetLink),
		TextBody: fmt.Sprintf("Password Reset\n\nClick the following link to reset your password:\n%s", p.ResetLink),
	}

	if err := h.emailService.Send(ctx, emailMsg); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	fmt.Printf("Sent password reset email to %s (%s)\n", user.Name, user.Email)

	return nil
}

// HandleExportUserData handles the user data export task
func (h *UserWorkerHandler) HandleExportUserData(ctx context.Context, payload worker.TaskPayload) error {
	var p ExportUserDataPayload

	data, _ := json.Marshal(payload)
	if err := json.Unmarshal(data, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if p.UserID == "" || p.Format == "" {
		return fmt.Errorf("missing required fields in payload")
	}

	// Get user
	user, err := h.userRepository.GetByID(ctx, p.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found: %s", p.UserID)
	}

	// In production, you would:
	// 1. Export user data in the specified format
	// 2. Store it in the destination (S3, local storage, etc.)
	// 3. Send a notification to the user with a download link

	fmt.Printf("Exporting user data for %s in format %s\n", p.UserID, p.Format)

	return nil
}

// HandleGenerateUserReport handles the user report generation task
func (h *UserWorkerHandler) HandleGenerateUserReport(ctx context.Context, payload worker.TaskPayload) error {
	var p GenerateUserReportPayload

	data, _ := json.Marshal(payload)
	if err := json.Unmarshal(data, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if p.UserID == "" || p.ReportType == "" {
		return fmt.Errorf("missing required fields in payload")
	}

	// Get user
	user, err := h.userRepository.GetByID(ctx, p.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found: %s", p.UserID)
	}

	// In production, you would:
	// 1. Gather data for the report based on report type
	// 2. Generate the report
	// 3. Store it in the destination
	// 4. Send notification to the user

	fmt.Printf("Generated %s report for user %s\n", p.ReportType, p.UserID)

	return nil
}
