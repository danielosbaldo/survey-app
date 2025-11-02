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
	RenderTemplate(c, "dashboard_section.gohtml", gin.H{
		"Title": "Dashboard",
	})
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
