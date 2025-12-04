package template

import (
	"fmt"
	"strings"
	"sync"
	"text/template"
)

// EmailTemplate represents a cached email template
type EmailTemplate struct {
	Subject      *template.Template
	HTMLBody     *template.Template
	TextBody     *template.Template
	Name         string
	RequiredKeys []string // Required data keys for validation
}

// TemplateLoader manages email templates with caching
type TemplateLoader struct {
	mu        sync.RWMutex
	templates map[string]*EmailTemplate
}

// NewTemplateLoader creates a new template loader
func NewTemplateLoader() *TemplateLoader {
	return &TemplateLoader{
		templates: make(map[string]*EmailTemplate),
	}
}

// RegisterTemplate registers a template with subject, HTML, and text versions
func (tl *TemplateLoader) RegisterTemplate(name string, subject, htmlBody, textBody string, requiredKeys []string) error {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	// Parse subject template
	subjectTmpl, err := template.New(name + "_subject").Parse(subject)
	if err != nil {
		return fmt.Errorf("failed to parse subject template: %w", err)
	}

	// Parse HTML template
	htmlTmpl, err := template.New(name + "_html").Parse(htmlBody)
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Parse text template
	textTmpl, err := template.New(name + "_text").Parse(textBody)
	if err != nil {
		return fmt.Errorf("failed to parse text template: %w", err)
	}

	tl.templates[name] = &EmailTemplate{
		Subject:      subjectTmpl,
		HTMLBody:     htmlTmpl,
		TextBody:     textTmpl,
		Name:         name,
		RequiredKeys: requiredKeys,
	}

	return nil
}

// GetTemplate retrieves a cached template
func (tl *TemplateLoader) GetTemplate(name string) (*EmailTemplate, error) {
	tl.mu.RLock()
	defer tl.mu.RUnlock()

	tmpl, ok := tl.templates[name]
	if !ok {
		return nil, fmt.Errorf("template not found: %s", name)
	}

	return tmpl, nil
}

// RenderTemplate renders a template with the given data
func (tl *TemplateLoader) RenderTemplate(name string, data map[string]interface{}) (subject, htmlBody, textBody string, err error) {
	tmpl, err := tl.GetTemplate(name)
	if err != nil {
		return "", "", "", err
	}

	// Validate required keys
	for _, key := range tmpl.RequiredKeys {
		if _, ok := data[key]; !ok {
			return "", "", "", fmt.Errorf("missing required template data: %s", key)
		}
	}

	// Render subject
	var subjectBuf strings.Builder
	if err := tmpl.Subject.Execute(&subjectBuf, data); err != nil {
		return "", "", "", fmt.Errorf("failed to render subject: %w", err)
	}

	// Render HTML
	var htmlBuf strings.Builder
	if err := tmpl.HTMLBody.Execute(&htmlBuf, data); err != nil {
		return "", "", "", fmt.Errorf("failed to render HTML: %w", err)
	}

	// Render text
	var textBuf strings.Builder
	if err := tmpl.TextBody.Execute(&textBuf, data); err != nil {
		return "", "", "", fmt.Errorf("failed to render text: %w", err)
	}

	return subjectBuf.String(), htmlBuf.String(), textBuf.String(), nil
}

// ListTemplates returns all registered template names
func (tl *TemplateLoader) ListTemplates() []string {
	tl.mu.RLock()
	defer tl.mu.RUnlock()

	names := make([]string, 0, len(tl.templates))
	for name := range tl.templates {
		names = append(names, name)
	}
	return names
}
