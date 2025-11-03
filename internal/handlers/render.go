package handlers

import (
	"encoding/json"
	"html/template"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/danielosbaldo/survey-app/assets"
)

var templates *template.Template

func init() {
	funcMap := template.FuncMap{
		"divFloat": func(a, b interface{}) float64 {
			var aFloat, bFloat float64

			switch v := a.(type) {
			case int:
				aFloat = float64(v)
			case int64:
				aFloat = float64(v)
			case float64:
				aFloat = v
			default:
				return 0
			}

			switch v := b.(type) {
			case int:
				bFloat = float64(v)
			case int64:
				bFloat = float64(v)
			case float64:
				bFloat = v
			default:
				return 0
			}

			if bFloat == 0 {
				return 0
			}
			return aFloat / bFloat
		},
		"toJSON": func(v interface{}) template.JS {
			b, err := json.Marshal(v)
			if err != nil {
				return template.JS("{}")
			}
			return template.JS(b)
		},
	}

	templates = template.Must(template.New("").Funcs(funcMap).ParseFS(assets.WebFS, "web/templates/*.gohtml", "web/templates/**/*.gohtml"))
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
