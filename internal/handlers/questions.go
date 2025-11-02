package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/danielosbaldo/survey-app/internal/models"
)

type QuestionHandler struct {
	DB *gorm.DB
}

func (h *QuestionHandler) Create(c *gin.Context) {
	prompt := c.PostForm("prompt")
	t := c.PostForm("type")
	order, _ := strconv.Atoi(c.PostForm("order"))

	// Create the question
	question := models.Question{Prompt: prompt, Type: t, OrderNum: order}
	if err := h.DB.Create(&question).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create question"})
		return
	}

	// Handle choices if provided
	choiceOrders := c.PostFormArray("choice_orders[]")
	choiceValues := c.PostFormArray("choice_values[]")
	choiceLabels := c.PostFormArray("choice_labels[]")

	// Debug: Print what we received
	if len(choiceValues) > 0 {
		println("DEBUG: Received", len(choiceValues), "choices for question", question.ID)
		for i := 0; i < len(choiceValues); i++ {
			println("  Choice", i, ":", choiceValues[i], "=", choiceLabels[i])
		}
	} else {
		println("DEBUG: No choices received for question", question.ID)
	}

	// Create choices
	for i := 0; i < len(choiceValues); i++ {
		if i >= len(choiceLabels) {
			continue
		}

		choiceOrder, _ := strconv.Atoi(choiceOrders[i])
		newChoice := models.Choice{
			QuestionID: question.ID,
			Label:      choiceLabels[i],
			Value:      choiceValues[i],
			OrderNum:   choiceOrder,
		}
		h.DB.Create(&newChoice)
	}

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

	// Handle choices
	choiceIDs := c.PostFormArray("choice_ids[]")
	choiceOrders := c.PostFormArray("choice_orders[]")
	choiceValues := c.PostFormArray("choice_values[]")
	choiceLabels := c.PostFormArray("choice_labels[]")

	// Debug: Print what we received
	println("DEBUG UPDATE: Question ID:", qID)
	println("DEBUG UPDATE: Received", len(choiceIDs), "choice IDs")
	println("DEBUG UPDATE: Received", len(choiceValues), "choice values")
	println("DEBUG UPDATE: Received", len(choiceLabels), "choice labels")
	for i := 0; i < len(choiceIDs); i++ {
		if i < len(choiceValues) && i < len(choiceLabels) {
			println("  Choice", i, "ID:", choiceIDs[i], "Value:", choiceValues[i], "Label:", choiceLabels[i])
		}
	}

	// Track existing choice IDs to keep
	keepChoiceIDs := make(map[uint]bool)

	// Update or create choices
	for i := 0; i < len(choiceIDs); i++ {
		if i >= len(choiceValues) || i >= len(choiceLabels) {
			continue
		}

		choiceID, _ := strconv.Atoi(choiceIDs[i])
		choiceOrder, _ := strconv.Atoi(choiceOrders[i])

		if choiceID == 0 {
			// Create new choice
			newChoice := models.Choice{
				QuestionID: uint(qID),
				Label:      choiceLabels[i],
				Value:      choiceValues[i],
				OrderNum:   choiceOrder,
			}
			if err := h.DB.Create(&newChoice).Error; err != nil {
				println("  ERROR creating new choice:", err.Error())
			} else {
				println("  SUCCESS: Created new choice ID:", newChoice.ID)
				// Add the new choice ID to the keep list so it doesn't get deleted
				keepChoiceIDs[newChoice.ID] = true
			}
		} else {
			// Update existing choice
			h.DB.Model(&models.Choice{}).Where("id = ?", choiceID).Updates(models.Choice{
				Label:    choiceLabels[i],
				Value:    choiceValues[i],
				OrderNum: choiceOrder,
			})
			keepChoiceIDs[uint(choiceID)] = true
		}
	}

	// Delete choices that were removed
	if len(keepChoiceIDs) > 0 {
		var idsToKeep []uint
		for id := range keepChoiceIDs {
			idsToKeep = append(idsToKeep, id)
		}
		h.DB.Where("question_id = ? AND id NOT IN ?", qID, idsToKeep).Delete(&models.Choice{})
	} else {
		// Delete all choices if none were kept
		h.DB.Where("question_id = ?", qID).Delete(&models.Choice{})
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

func (h *QuestionHandler) Edit(c *gin.Context) {
	id := c.Param("id")
	qID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid question ID"})
		return
	}

	var question models.Question
	if err := h.DB.Preload("Choices", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_num ASC")
	}).First(&question, qID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Question not found"})
		return
	}

	RenderTemplate(c, "question_edit_row.gohtml", gin.H{
		"Question": question,
	})
}

func (h *QuestionHandler) Section(c *gin.Context) {
	RenderTemplate(c, "questions_section.gohtml", getAdminData(h.DB))
}

func renderQuestionsTable(c *gin.Context, db *gorm.DB) {
	RenderTemplate(c, "questions_table.gohtml", getAdminData(db))
}
