package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/example/heladeria/internal/models"
)

type QuestionHandler struct {
	DB *gorm.DB
}

func (h *QuestionHandler) Create(c *gin.Context) {
	prompt := c.PostForm("prompt")
	t := c.PostForm("type")
	order, _ := strconv.Atoi(c.PostForm("order"))
	h.DB.Create(&models.Question{Prompt: prompt, Type: t, OrderNum: order})

	renderQuestionsTable(c, h.DB)
}

func (h *QuestionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	qID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid question ID"})
		return
	}

	prompt := c.PostForm("prompt")
	qType := c.PostForm("type")
	order, _ := strconv.Atoi(c.PostForm("order"))

	question := models.Question{
		Prompt:   prompt,
		Type:     qType,
		OrderNum: order,
	}

	if err := h.DB.Model(&models.Question{}).Where("id = ?", qID).Updates(question).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update question"})
		return
	}

	renderQuestionsTable(c, h.DB)
}

func (h *QuestionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	qID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid question ID"})
		return
	}

	h.DB.Where("question_id = ?", qID).Delete(&models.Choice{})

	if err := h.DB.Delete(&models.Question{}, qID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete question"})
		return
	}

	renderQuestionsTable(c, h.DB)
}

func (h *QuestionHandler) CreateChoice(c *gin.Context) {
	qid, _ := strconv.Atoi(c.PostForm("question_id"))
	label := c.PostForm("label")
	value := c.PostForm("value")
	order, _ := strconv.Atoi(c.PostForm("order"))
	h.DB.Create(&models.Choice{QuestionID: uint(qid), Label: label, Value: value, OrderNum: order})

	RenderTemplate(c, "partials_admin_refresh.gohtml", getAdminData(h.DB))
}

func (h *QuestionHandler) Section(c *gin.Context) {
	RenderTemplate(c, "questions_section.gohtml", getAdminData(h.DB))
}

func renderQuestionsTable(c *gin.Context, db *gorm.DB) {
	RenderTemplate(c, "questions_table.gohtml", getAdminData(db))
}
