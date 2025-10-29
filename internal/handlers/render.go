package handlers

import (
	"html/template"

	"github.com/gin-gonic/gin"

	"github.com/example/heladeria/assets"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFS(assets.WebFS, "web/templates/*.gohtml", "web/templates/**/*.gohtml"))
}

func RenderTemplate(c *gin.Context, name string, data gin.H) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	templates.ExecuteTemplate(c.Writer, name, data)
}
