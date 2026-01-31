package service

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/gofund/shared/models"
)

// RenderService handles email template rendering
type RenderService interface {
	Render(emailType models.EmailType, data map[string]interface{}) (string, error)
}

type renderService struct {
	templateDir string
}

// NewRenderService creates a new render service
func NewRenderService(templateDir string) RenderService {
	return &renderService{
		templateDir: templateDir,
	}
}

// Render renders an email template with the given data
func (s *renderService) Render(emailType models.EmailType, data map[string]interface{}) (string, error) {
	// Add common data
	if _, ok := data["Year"]; !ok {
		data["Year"] = time.Now().Year()
	}

	layoutPath := filepath.Join(s.templateDir, "layout.html")
	templatePath := filepath.Join(s.templateDir, fmt.Sprintf("%s.html", emailType))

	tmpl, err := template.ParseFiles(layoutPath, templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse templates: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "layout.html", data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
