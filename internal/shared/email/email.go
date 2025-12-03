package email

import "context"

// EmailType defines the type of email to send
type EmailType string

const (
	EmailTypeWelcome       EmailType = "welcome"
	EmailTypePasswordReset EmailType = "password_reset"
	EmailTypeVerification  EmailType = "verification"
	EmailTypeNotification  EmailType = "notification"
)

// Email represents an email message
type Email struct {
	To           []string               `json:"to"`
	CC           []string               `json:"cc"`
	BCC          []string               `json:"bcc"`
	From         string                 `json:"from"`
	ReplyTo      string                 `json:"reply_to"`
	Subject      string                 `json:"subject"`
	TextBody     string                 `json:"text_body"`
	HTMLBody     string                 `json:"html_body"`
	Headers      map[string]string      `json:"headers"`
	Attachments  []Attachment           `json:"attachments"`
	TemplateData map[string]interface{} `json:"template_data"`
}

// Attachment represents an email attachment
type Attachment struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
	MimeType string `json:"mime_type"`
}

// EmailService defines the interface for sending emails
type EmailService interface {
	// Send sends a single email
	Send(ctx context.Context, email *Email) error

	// SendBatch sends multiple emails in batch
	SendBatch(ctx context.Context, emails []*Email) error

	// SendTemplate sends an email using a template
	SendTemplate(ctx context.Context, to []string, templateID string, data map[string]interface{}) error

	// ValidateEmail checks if an email address is valid
	ValidateEmail(email string) error

	// Health checks if the email service is healthy
	Health(ctx context.Context) error
}

// SendResult represents the result of sending an email
type SendResult struct {
	EmailID   string
	MessageID string
	Status    string
	Error     string
}
