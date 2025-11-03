package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/danielosbaldo/survey-app/internal/models"
)

type EmployeeHandler struct {
	DB *gorm.DB
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.String(400, "nombre requerido")
		return
	}
	active := c.PostForm("active") == "true"
	shopIDs := c.PostFormArray("shop_ids")

	emp := models.Employee{Name: name, Active: active}
	h.DB.Create(&emp)

	var shops []models.Shop
	for _, idStr := range shopIDs {
		if id, err := strconv.Atoi(idStr); err == nil {
			var shop models.Shop
			if err := h.DB.First(&shop, id).Error; err == nil {
				shops = append(shops, shop)
			}
		}
	}
	if len(shops) > 0 {
		h.DB.Model(&emp).Association("Shops").Replace(shops)
	}

	renderEmployeesTable(c, h.DB)
}

func (h *EmployeeHandler) Edit(c *gin.Context) {
	id := c.Param("id")
	empID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid employee ID"})
		return
	}

	var emp models.Employee
	if err := h.DB.Preload("Shops").Preload("Shops.Ciudad").First(&emp, empID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Employee not found"})
		return
	}

	var ciudades []models.Ciudad
	h.DB.Find(&ciudades)

	// Get the ciudad_id from the first shop (if exists)
	var selectedCiudadID uint
	var availableShops []models.Shop
	if len(emp.Shops) > 0 {
		selectedCiudadID = emp.Shops[0].CiudadID
		h.DB.Where("ciudad_id = ?", selectedCiudadID).Find(&availableShops)
	} else {
		// If no shops assigned, still need empty array
		availableShops = []models.Shop{}
	}

	// Create a map of selected shop IDs for easy lookup
	selectedShopIDs := make(map[uint]bool)
	for _, shop := range emp.Shops {
		selectedShopIDs[shop.ID] = true
	}

	// Debug logging
	c.Writer.Header().Set("X-Debug-Employee-ID", id)
	c.Writer.Header().Set("X-Debug-Shops-Count", strconv.Itoa(len(emp.Shops)))
	c.Writer.Header().Set("X-Debug-Selected-Ciudad", strconv.Itoa(int(selectedCiudadID)))

	RenderTemplate(c, "employee_edit_row.gohtml", gin.H{
		"Employee":         emp,
		"Ciudades":         ciudades,
		"SelectedCiudadID": selectedCiudadID,
		"AvailableShops":   availableShops,
		"IsShopSelected": func(shopID uint) bool {
			return selectedShopIDs[shopID]
		},
	})
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	id := c.Param("id")
	empID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid employee ID"})
		return
	}

	var emp models.Employee
	if err := h.DB.First(&emp, empID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Employee not found"})
		return
	}

	name := c.PostForm("name")
	active := c.PostForm("active") == "true"
	shopIDs := c.PostFormArray("shop_ids")

	// Update employee basic info
	emp.Name = name
	emp.Active = active

	if err := h.DB.Save(&emp).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update employee"})
		return
	}

	// Update shop associations - clear all and add new ones
	h.DB.Model(&emp).Association("Shops").Clear()

	if len(shopIDs) > 0 {
		var shops []models.Shop
		for _, shopIDStr := range shopIDs {
			shopID, _ := strconv.Atoi(shopIDStr)
			var shop models.Shop
			if h.DB.First(&shop, shopID).Error == nil {
				shops = append(shops, shop)
			}
		}
		if len(shops) > 0 {
			h.DB.Model(&emp).Association("Shops").Append(shops)
		}
	}

	renderEmployeesTable(c, h.DB)
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	empID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid employee ID"})
		return
	}

	if err := h.DB.Delete(&models.Employee{}, empID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete employee"})
		return
	}

	renderEmployeesTable(c, h.DB)
}

func (h *EmployeeHandler) Toggle(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	var emp models.Employee
	if err := h.DB.First(&emp, id).Error; err != nil {
		c.String(404, "empleado no encontrado")
		return
	}
	emp.Active = !emp.Active
	h.DB.Save(&emp)

	RenderTemplate(c, "partials_admin_refresh.gohtml", getAdminData(h.DB))
}

func (h *EmployeeHandler) Section(c *gin.Context) {
	RenderTemplate(c, "employees_section.gohtml", getAdminData(h.DB))
}

func renderEmployeesTable(c *gin.Context, db *gorm.DB) {
	RenderTemplate(c, "employees_content.gohtml", getAdminData(db))
}
