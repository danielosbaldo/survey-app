package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/example/heladeria/internal/models"
)

type SurveyHandler struct {
	DB *gorm.DB
}

func (h *SurveyHandler) GetForm(c *gin.Context) {
	slug := c.Param("slug")
	ciudad := c.Param("ciudad")

	var shop models.Shop
	if err := h.DB.Joins("JOIN ciudades ON shops.ciudad_id = ciudades.id").
		Where("LOWER(ciudades.nombre) = LOWER(?) AND shops.slug = ?", ciudad, slug).
		First(&shop).Error; err != nil {
		c.String(404, "Sucursal no encontrada")
		return
	}

	var emps []models.Employee
	h.DB.Joins("JOIN employee_shops es ON es.employee_id = employees.id").
		Where("es.shop_id = ? AND employees.active = ?", shop.ID, true).
		Order("employees.name asc").
		Find(&emps)

	var questions []models.Question
	h.DB.Preload("Choices").Order("order_num asc").Find(&questions)

	errorMsg := ""
	errorType := c.Query("error")
	switch errorType {
	case "employee_required":
		errorMsg = "Por favor, selecciona un empleado"
	case "employee_not_found":
		errorMsg = "Empleado no encontrado. Por favor, selecciona un empleado de la lista"
	}

	preservedData := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if key != "error" && len(values) > 0 {
			preservedData[key] = values[0]
		}
	}

	RenderTemplate(c, "form.gohtml", gin.H{
		"Title":         "Encuesta",
		"Shop":          shop,
		"Ciudad":        ciudad,
		"Employees":     emps,
		"Questions":     questions,
		"Error":         errorMsg,
		"PreservedData": preservedData,
	})
}

func (h *SurveyHandler) PostForm(c *gin.Context) {
	slug := c.Param("slug")
	ciudad := c.Param("ciudad")

	var shop models.Shop
	if err := h.DB.Joins("JOIN ciudads ON shops.ciudad_id = ciudads.id").
		Where("LOWER(ciudads.nombre) = LOWER(?) AND shops.slug = ?", ciudad, slug).
		First(&shop).Error; err != nil {
		c.String(404, "Sucursal no encontrada")
		return
	}

	empName := c.PostForm("employee_name")
	if empName == "" {
		c.Redirect(302, fmt.Sprintf("/sucursal/%s/%s/encuesta?error=employee_required&%s", ciudad, slug, c.Request.URL.RawQuery))
		return
	}

	var emp models.Employee
	if err := h.DB.Joins("JOIN employee_shops es ON es.employee_id = employees.id").
		Where("es.shop_id = ? AND employees.name = ? AND employees.active = ?", shop.ID, empName, true).
		First(&emp).Error; err != nil {
		formData := c.Request.Form.Encode()
		c.Redirect(302, fmt.Sprintf("/sucursal/%s/%s/encuesta?error=employee_not_found&%s", ciudad, slug, formData))
		return
	}

	answers := models.JSONB{}
	for k, v := range c.Request.PostForm {
		if k == "employee_name" {
			continue
		}
		if len(v) > 0 {
			answers[k] = v[0]
		}
	}

	h.DB.Create(&models.Response{ShopID: shop.ID, EmployeeID: emp.ID, Answers: answers, UserAgent: c.Request.UserAgent()})
	c.Data(200, "text/html; charset=utf-8", []byte(`<div style="font-family:system-ui;padding:24px;text-align:center"><h2>¡Gracias por tu opinión! 🍨</h2><a href="">Enviar otra</a></div>`))
}
