package handlers

import (
	"html/template"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/danielosbaldo/survey-app/assets"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFS(assets.WebFS, "web/templates/*.gohtml", "web/templates/**/*.gohtml"))
}

// GetAppName returns the configured application name from environment or default
func GetAppName() string {
	if name := os.Getenv("APP_NAME"); name != "" {
		return name
	}
	return "Survey App"
}

func RenderTemplate(c *gin.Context, name string, data gin.H) {
	// Ensure AppName is always available in templates
	if data == nil {
		data = gin.H{}
	}
	if _, exists := data["AppName"]; !exists {
		data["AppName"] = GetAppName()
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(c.Writer, name, data)
}
