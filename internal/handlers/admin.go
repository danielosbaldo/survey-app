package handlers

import (
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/danielosbaldo/survey-app/internal/models"
)

type AdminHandler struct {
	DB *gorm.DB
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	var total int64
	h.DB.Model(&models.Response{}).Count(&total)

	avg := 0.0
	var responses []models.Response
	h.DB.Find(&responses)
	sum, cnt := 0, 0
	for _, r := range responses {
		if v, ok := r.Answers["amabilidad"]; ok {
			if str, ok2 := v.(string); ok2 {
				if iv, err := strconv.Atoi(str); err == nil {
					sum += iv
					cnt++
				}
			}
		}
	}
	if cnt > 0 {
		avg = float64(sum) / float64(cnt)
	}

	data := getAdminData(h.DB)
	data["Title"] = "Admin"
	data["Total"] = total
	data["AvgAmabilidad"] = avg

	RenderTemplate(c, "admin.gohtml", data)
}

func (h *AdminHandler) DashboardSection(c *gin.Context) {
	shopIDStr := c.Query("shop_id")
	var shopID uint
	if shopIDStr != "" {
		if id, err := strconv.Atoi(shopIDStr); err == nil {
			shopID = uint(id)
		}
	}

	data := getDashboardData(h.DB, shopID)
	data["Title"] = "Dashboard"
	data["SelectedShopID"] = shopID

	// Get all shops for the filter dropdown
	var shops []models.Shop
	h.DB.Preload("Ciudad").Order("name").Find(&shops)
	data["Shops"] = shops

	RenderTemplate(c, "dashboard_section.gohtml", data)
}

func (h *AdminHandler) PartialTable(c *gin.Context) {
	var rows []models.Response
	h.DB.Order("created_at desc").Limit(50).Find(&rows)
	RenderTemplate(c, "partials_table.gohtml", gin.H{"Rows": rows})
}

func (h *AdminHandler) PartialKPIs(c *gin.Context) {
	var total int64
	h.DB.Model(&models.Response{}).Count(&total)

	avg := 0.0
	var responses []models.Response
	h.DB.Find(&responses)
	sum, cnt := 0, 0
	for _, r := range responses {
		if v, ok := r.Answers["amabilidad"]; ok {
			if str, ok2 := v.(string); ok2 {
				if iv, err := strconv.Atoi(str); err == nil {
					sum += iv
					cnt++
				}
			}
		}
	}
	if cnt > 0 {
		avg = float64(sum) / float64(cnt)
	}

	RenderTemplate(c, "partials_kpis.gohtml", gin.H{"Total": total, "AvgAmabilidad": avg})
}

func getAdminData(db *gorm.DB) gin.H {
	var rows []models.Response
	db.Order("created_at desc").Limit(50).Find(&rows)
	var shops []models.Shop
	db.Preload("Ciudad").Find(&shops)
	var employees []models.Employee
	db.Preload("Shops").Find(&employees)
	sort.Slice(employees, func(i, j int) bool { return employees[i].Name < employees[j].Name })
	var ciudades []models.Ciudad
	db.Find(&ciudades)
	var questions []models.Question
	db.Preload("Choices", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_num ASC")
	}).Order("order_num asc").Find(&questions)

	// Debug: Print how many choices each question has
	for _, q := range questions {
		println("DEBUG getAdminData: Question", q.ID, "has", len(q.Choices), "choices")
	}

	return gin.H{
		"Rows":      rows,
		"Shops":     shops,
		"Employees": employees,
		"Questions": questions,
		"Ciudades":  ciudades,
	}
}

func getDashboardData(db *gorm.DB, shopID uint) gin.H {
	// Build query with optional shop filter
	responseQuery := db.Model(&models.Response{})
	if shopID > 0 {
		responseQuery = responseQuery.Where("shop_id = ?", shopID)
	}

	// Total responses
	var totalResponses int64
	responseQuery.Count(&totalResponses)

	// Total employees (active)
	var totalEmployees int64
	employeeQuery := db.Model(&models.Employee{}).Where("active = ?", true)
	if shopID > 0 {
		// Count employees assigned to this shop
		employeeQuery = employeeQuery.Joins("JOIN employee_shops ON employee_shops.employee_id = employees.id").
			Where("employee_shops.shop_id = ?", shopID)
	}
	employeeQuery.Count(&totalEmployees)

	// Total shops
	var totalShops int64
	if shopID > 0 {
		totalShops = 1 // Just the selected shop
	} else {
		db.Model(&models.Shop{}).Count(&totalShops)
	}

	// Get responses for analysis
	var responses []models.Response
	query := db.Preload("Shop").Preload("Shop.Ciudad").Order("created_at desc")
	if shopID > 0 {
		query = query.Where("shop_id = ?", shopID)
	}
	query.Find(&responses)

	// Get all questions for analysis
	var questions []models.Question
	db.Preload("Choices", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_num ASC")
	}).Order("order_num asc").Find(&questions)

	// Responses by shop
	responsesByShop := make(map[string]int)
	for _, r := range responses {
		if r.Shop.ID != 0 {
			responsesByShop[r.Shop.Name]++
		}
	}

	// Responses over time (last 7 days)
	responsesOverTime := make(map[string]int)
	for _, r := range responses {
		dateKey := r.CreatedAt.Format("2006-01-02")
		responsesOverTime[dateKey]++
	}

	// Question response stats - use array to preserve order and include IDs
	type QuestionStat struct {
		ID     uint              `json:"id"`
		Prompt string            `json:"prompt"`
		Stats  map[string]int    `json:"stats"`
	}
	questionStats := []QuestionStat{}

	for _, q := range questions {
		stats := make(map[string]int)
		for _, r := range responses {
			// Try both "qN" format and just "N" format
			key := "q" + strconv.Itoa(int(q.ID))
			val, ok := r.Answers[key]
			if !ok {
				// Try without "q" prefix
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
				}
			}
		}

		questionStats = append(questionStats, QuestionStat{
			ID:     q.ID,
			Prompt: q.Prompt,
			Stats:  stats,
		})
	}

	// Recent responses (last 10)
	var recentResponses []models.Response
	db.Preload("Shop").Preload("Shop.Ciudad").Order("created_at desc").Limit(10).Find(&recentResponses)

	return gin.H{
		"TotalResponses":    totalResponses,
		"TotalEmployees":    totalEmployees,
		"TotalShops":        totalShops,
		"ResponsesByShop":   responsesByShop,
		"ResponsesOverTime": responsesOverTime,
		"QuestionStats":     questionStats,
		"RecentResponses":   recentResponses,
		"Questions":         questions,
	}
}
