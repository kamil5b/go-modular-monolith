package email

import (
	"context"
	"fmt"
	"regexp"
)

// NoOpEmailService is a no-op implementation that doesn't send emails
// Used for testing and development without a real email provider
type NoOpEmailService struct{}

// NewNoOpEmailService creates a new no-op email service
func NewNoOpEmailService() *NoOpEmailService {
	return &NoOpEmailService{}
}

// Send is a no-op
func (s *NoOpEmailService) Send(ctx context.Context, email *Email) error {
	if email == nil {
		return fmt.Errorf("email cannot be nil")
	}
	fmt.Printf("[NoOp] Would send email to %v with subject: %s\n", email.To, email.Subject)
	return nil
}

// SendBatch is a no-op
func (s *NoOpEmailService) SendBatch(ctx context.Context, emails []*Email) error {
	if emails == nil {
		return fmt.Errorf("emails cannot be nil")
	}
	fmt.Printf("[NoOp] Would send %d emails in batch\n", len(emails))
	return nil
}

// SendTemplate is a no-op
func (s *NoOpEmailService) SendTemplate(ctx context.Context, to []string, templateID string, data map[string]interface{}) error {
	if len(to) == 0 {
		return fmt.Errorf("to cannot be empty")
	}
	fmt.Printf("[NoOp] Would send template email to %v (template: %s)\n", to, templateID)
	return nil
}

// ValidateEmail validates an email address format
func (s *NoOpEmailService) ValidateEmail(email string) error {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}
	return nil
}

// Health is always healthy for no-op
func (s *NoOpEmailService) Health(ctx context.Context) error {
	return nil
}
