package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/danielosbaldo/survey-app/assets"
	"github.com/danielosbaldo/survey-app/internal/handlers"
)

type Server struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Server {
	return &Server{DB: db}
}

func (s *Server) Router() *gin.Engine {
	g := gin.Default()

	// Static assets
	g.HEAD("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/") })
	g.StaticFS("/assets", http.FS(assets.WebFS))

	// Health check endpoint
	g.GET("/health", func(c *gin.Context) {
		// Check database connection
		sqlDB, err := s.DB.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": "database connection error"})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": "database ping failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Home
	g.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `<div style="font-family:system-ui;padding:48px;text-align:center"><h1>Bienvenido</h1></div>`)
	})

	// Survey handlers
	survey := &handlers.SurveyHandler{DB: s.DB}
	g.GET("/sucursal/*path", survey.GetFormRouter)
	g.POST("/sucursal/:ciudad/:slug/encuesta", survey.PostForm)

	// Admin handlers
	admin := g.Group("/admin")
	{
		adminHandler := &handlers.AdminHandler{DB: s.DB}
		shopHandler := &handlers.ShopHandler{DB: s.DB}
		employeeHandler := &handlers.EmployeeHandler{DB: s.DB}
		questionHandler := &handlers.QuestionHandler{DB: s.DB}

		// Dashboard
		admin.GET("", adminHandler.Dashboard)
		admin.GET("/dashboard-section", adminHandler.DashboardSection)
		admin.GET("/table", adminHandler.PartialTable)
		admin.GET("/kpis", adminHandler.PartialKPIs)

		// Shops
		admin.GET("/shops-section", shopHandler.Section)
		admin.POST("/shops", shopHandler.Create)
		admin.PUT("/shops/:id", shopHandler.Update)
		admin.DELETE("/shops/:id", shopHandler.Delete)
		admin.GET("/shops-by-ciudad", shopHandler.GetByCiudad)

		// Employees
		admin.GET("/employees-section", employeeHandler.Section)
		admin.GET("/employees/:id/edit", employeeHandler.Edit)
		admin.POST("/employees", employeeHandler.Create)
		admin.PUT("/employees/:id", employeeHandler.Update)
		admin.DELETE("/employees/:id", employeeHandler.Delete)
		admin.POST("/toggle-employee", employeeHandler.Toggle)

		// Questions
		admin.GET("/questions-section", questionHandler.Section)
		admin.GET("/questions/:id/edit", questionHandler.Edit)
		admin.POST("/questions", questionHandler.Create)
		admin.PUT("/questions/:id", questionHandler.Update)
		admin.DELETE("/questions/:id", questionHandler.Delete)
		admin.POST("/choices", questionHandler.CreateChoice)
	}

	return g
}
