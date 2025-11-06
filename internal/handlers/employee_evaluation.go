package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/danielosbaldo/survey-app/internal/models"
)

type EmployeeEvaluationHandler struct {
	DB *gorm.DB
}

func (h *EmployeeEvaluationHandler) Section(c *gin.Context) {
	shopIDStr := c.Query("shop_id")
	employeeIDStr := c.Query("employee_id")

	var shopID, employeeID uint
	if shopIDStr != "" {
		if id, err := strconv.Atoi(shopIDStr); err == nil {
			shopID = uint(id)
		}
	}
	if employeeIDStr != "" {
		if id, err := strconv.Atoi(employeeIDStr); err == nil {
			employeeID = uint(id)
		}
	}

	data := h.getEvaluationData(shopID, employeeID)
	data["Title"] = "Evaluación de Empleados"
	data["SelectedShopID"] = shopID
	data["SelectedEmployeeID"] = employeeID

	// Get all shops for the filter dropdown
	var shops []models.Shop
	h.DB.Preload("Ciudad").Order("name").Find(&shops)
	data["Shops"] = shops

	// Get employees based on shop filter
	var employees []models.Employee
	employeeQuery := h.DB.Where("active = ?", true)
	if shopID > 0 {
		employeeQuery = employeeQuery.Joins("JOIN employee_shops ON employee_shops.employee_id = employees.id").
			Where("employee_shops.shop_id = ?", shopID)
	}
	employeeQuery.Order("name").Find(&employees)
	data["Employees"] = employees

	RenderTemplate(c, "employee_evaluation_section.gohtml", data)
}

func (h *EmployeeEvaluationHandler) getEvaluationData(shopID, employeeID uint) gin.H {
	// Build query with optional filters
	responseQuery := h.DB.Model(&models.Response{})
	if shopID > 0 {
		responseQuery = responseQuery.Where("shop_id = ?", shopID)
	}
	if employeeID > 0 {
		responseQuery = responseQuery.Where("employee_id = ?", employeeID)
	}

	// Total responses for the employee
	var totalResponses int64
	responseQuery.Count(&totalResponses)

	// Get responses for analysis
	var responses []models.Response
	query := h.DB.Preload("Shop").Preload("Shop.Ciudad").Preload("Employee").Order("created_at desc")
	if shopID > 0 {
		query = query.Where("shop_id = ?", shopID)
	}
	if employeeID > 0 {
		query = query.Where("employee_id = ?", employeeID)
	}
	query.Find(&responses)

	// Get all questions for analysis
	var questions []models.Question
	h.DB.Preload("Choices", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_num ASC")
	}).Order("order_num asc").Find(&questions)

	// Calculate average scores per question
	type QuestionAverage struct {
		ID      uint               `json:"id"`
		Prompt  string             `json:"prompt"`
		Average float64            `json:"average"`
		Total   int                `json:"total"`
		Stats   map[string]int     `json:"stats"`
	}

	questionAverages := []QuestionAverage{}

	for _, q := range questions {
		stats := make(map[string]int)
		sum := 0.0
		count := 0

		for _, r := range responses {
			// Try both "qN" format and just "N" format
			key := "q" + strconv.Itoa(int(q.ID))
			val, ok := r.Answers[key]
			if !ok {
				val, ok = r.Answers[strconv.Itoa(int(q.ID))]
			}

			if ok {
				valStr := ""
				switch v := val.(type) {
				case string:
					valStr = v
				case float64:
					valStr = strconv.Itoa(int(v))
				}
				if valStr != "" {
					stats[valStr]++
					// Try to parse as number for average calculation
					if numVal, err := strconv.ParseFloat(valStr, 64); err == nil {
						sum += numVal
						count++
					}
				}
			}
		}

		avg := 0.0
		if count > 0 {
			avg = sum / float64(count)
		}

		questionAverages = append(questionAverages, QuestionAverage{
			ID:      q.ID,
			Prompt:  q.Prompt,
			Average: avg,
			Total:   count,
			Stats:   stats,
		})
	}

	// Responses over time
	responsesOverTime := make(map[string]int)
	for _, r := range responses {
		dateKey := r.CreatedAt.Format("2006-01-02")
		responsesOverTime[dateKey]++
	}

	// Recent responses for the employee
	var recentResponses []models.Response
	recentQuery := h.DB.Preload("Shop").Preload("Shop.Ciudad").Preload("Employee").Order("created_at desc").Limit(10)
	if shopID > 0 {
		recentQuery = recentQuery.Where("shop_id = ?", shopID)
	}
	if employeeID > 0 {
		recentQuery = recentQuery.Where("employee_id = ?", employeeID)
	}
	recentQuery.Find(&recentResponses)

	// Get employee info if selected
	var selectedEmployee *models.Employee
	if employeeID > 0 {
		var emp models.Employee
		if err := h.DB.Preload("Shops").First(&emp, employeeID).Error; err == nil {
			selectedEmployee = &emp
		}
	}

	return gin.H{
		"TotalResponses":     totalResponses,
		"ResponsesOverTime":  responsesOverTime,
		"QuestionAverages":   questionAverages,
		"RecentResponses":    recentResponses,
		"Questions":          questions,
		"SelectedEmployee":   selectedEmployee,
	}
}
