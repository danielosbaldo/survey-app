package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/danielosbaldo/survey-app/internal/models"
)

type ShopHandler struct {
	DB *gorm.DB
}

func (h *ShopHandler) Create(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.String(400, "nombre requerido")
		return
	}
	slug := c.PostForm("slug")
	if slug == "" {
		slug = name
	}
	ciudadID, _ := strconv.Atoi(c.PostForm("ciudad_id"))

	var existingShop models.Shop
	if err := h.DB.Where("slug = ?", slug).First(&existingShop).Error; err == nil {
		c.String(400, "El slug '"+slug+"' ya existe. Por favor usa uno diferente.")
		return
	}

	shop := models.Shop{Name: name, Slug: slug, CiudadID: uint(ciudadID)}
	if err := h.DB.Create(&shop).Error; err != nil {
		c.String(500, "Error al crear la sucursal: "+err.Error())
		return
	}

	renderShopsTable(c, h.DB)
}

func (h *ShopHandler) Update(c *gin.Context) {
	id := c.Param("id")
	shopID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid shop ID"})
		return
	}

	var currentShop models.Shop
	if err := h.DB.First(&currentShop, shopID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Shop not found"})
		return
	}

	name := c.PostForm("name")
	slug := c.PostForm("slug")
	ciudadID, _ := strconv.Atoi(c.PostForm("ciudad_id"))

	if slug != currentShop.Slug {
		var existingShop models.Shop
		if err := h.DB.Where("slug = ?", slug).First(&existingShop).Error; err == nil {
			c.JSON(400, gin.H{"error": "El slug '" + slug + "' ya existe. Por favor usa uno diferente."})
			return
		}
	}

	updates := make(map[string]interface{})
	if name != currentShop.Name {
		updates["name"] = name
	}
	if slug != currentShop.Slug {
		updates["slug"] = slug
	}
	if uint(ciudadID) != currentShop.CiudadID {
		updates["ciudad_id"] = uint(ciudadID)
	}

	if len(updates) > 0 {
		if err := h.DB.Model(&currentShop).Updates(updates).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update shop: " + err.Error()})
			return
		}
	}

	renderShopsTable(c, h.DB)
}

func (h *ShopHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	shopID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid shop ID"})
		return
	}

	if err := h.DB.Delete(&models.Shop{}, shopID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete shop"})
		return
	}

	renderShopsTable(c, h.DB)
}

func (h *ShopHandler) GetByCiudad(c *gin.Context) {
	ciudadID := c.Query("ciudad_id")
	if ciudadID == "" {
		ciudadID = c.PostForm("ciudad_id")
	}

	// Check if this is for an edit form - check both query and JSON body
	employeeID := c.Query("employee_id")
	if employeeID == "" {
		// Try to get from hx-vals JSON
		var vals map[string]string
		c.ShouldBindJSON(&vals)
		employeeID = vals["employee_id"]
	}

	targetID := "employee-shops-dropdown-container"
	if employeeID != "" {
		targetID = "edit-shops-dropdown-" + employeeID
	}

	var data gin.H
	if ciudadID == "" {
		data = gin.H{
			"Shops":    nil,
			"Message":  "Primero selecciona una ciudad",
			"TargetID": targetID,
		}
	} else {
		var shops []models.Shop
		h.DB.Where("ciudad_id = ?", ciudadID).Find(&shops)

		if len(shops) == 0 {
			data = gin.H{
				"Shops":    nil,
				"Message":  "No hay sucursales en esta ciudad",
				"TargetID": targetID,
			}
		} else {
			data = gin.H{
				"Shops":    shops,
				"Message":  "",
				"TargetID": targetID,
			}
		}
	}

	// Use different template for edit vs create
	if employeeID != "" {
		RenderTemplate(c, "shops_dropdown_edit.gohtml", data)
	} else {
		RenderTemplate(c, "shops_dropdown.gohtml", data)
	}
}

func (h *ShopHandler) Section(c *gin.Context) {
	var shops []models.Shop
	h.DB.Preload("Ciudad").Find(&shops)
	var ciudades []models.Ciudad
	h.DB.Find(&ciudades)

	RenderTemplate(c, "shops_section.gohtml", gin.H{
		"Title":    "Sucursales",
		"Shops":    shops,
		"Ciudades": ciudades,
	})
}

func renderShopsTable(c *gin.Context, db *gorm.DB) {
	var shops []models.Shop
	db.Preload("Ciudad").Find(&shops)
	var ciudades []models.Ciudad
	db.Find(&ciudades)

	RenderTemplate(c, "shops_table.gohtml", gin.H{
		"Shops":    shops,
		"Ciudades": ciudades,
	})
}
