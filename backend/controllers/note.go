package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"collaborative-notes/config"
	"collaborative-notes/models"
)

func GetNotes(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var notes []models.Note
	// Get notes owned by user or collaborated on
	query := config.DB.Where("owner_id = ?", userID).
		Or("id IN (SELECT note_id FROM collaborations WHERE user_id = ? AND deleted_at IS NULL)", userID).
		Preload("Owner").
		Preload("Collaborations.User")
	
	if err := query.Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	c.JSON(http.StatusOK, notes)
}

func CreateNote(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req models.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := models.Note{
		Title:    req.Title,
		Content:  req.Content,
		OwnerID:  userID,
		IsPublic: req.IsPublic,
		Version:  1,
	}

	if err := config.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	// Load relationships
	config.DB.Preload("Owner").Preload("Collaborations.User").First(&note, note.ID)

	c.JSON(http.StatusCreated, note)
}

func GetNote(c *gin.Context) {
	userID := c.GetUint("user_id")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var note models.Note
	query := config.DB.Where("id = ?", noteID).
		Where("owner_id = ? OR is_public = true OR id IN (SELECT note_id FROM collaborations WHERE user_id = ? AND deleted_at IS NULL)", userID, userID).
		Preload("Owner").
		Preload("Collaborations.User")
	
	if err := query.First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func UpdateNote(c *gin.Context) {
	userID := c.GetUint("user_id")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var req models.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var note models.Note
	// Check if user has write permission
	query := config.DB.Where("id = ?", noteID).
		Where("owner_id = ? OR id IN (SELECT note_id FROM collaborations WHERE user_id = ? AND permission = 'write' AND deleted_at IS NULL)", userID, userID)
	
	if err := query.First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or insufficient permissions"})
		return
	}

	// Update fields if provided
	if req.Title != nil {
		note.Title = *req.Title
	}
	if req.Content != nil {
		note.Content = *req.Content
		note.Version++
	}
	if req.IsPublic != nil && note.OwnerID == userID {
		note.IsPublic = *req.IsPublic
	}

	if err := config.DB.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	// Load relationships
	config.DB.Preload("Owner").Preload("Collaborations.User").First(&note, note.ID)

	c.JSON(http.StatusOK, note)
}

func DeleteNote(c *gin.Context) {
	userID := c.GetUint("user_id")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var note models.Note
	if err := config.DB.Where("id = ? AND owner_id = ?", noteID, userID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	if err := config.DB.Delete(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

func ShareNote(c *gin.Context) {
	userID := c.GetUint("user_id")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var req models.ShareNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user owns the note
	var note models.Note
	if err := config.DB.Where("id = ? AND owner_id = ?", noteID, userID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	// Find user to share with
	var targetUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&targetUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if collaboration already exists
	var existingCollab models.Collaboration
	if err := config.DB.Where("note_id = ? AND user_id = ?", noteID, targetUser.ID).First(&existingCollab).Error; err == nil {
		// Update existing collaboration
		existingCollab.Permission = req.Permission
		config.DB.Save(&existingCollab)
	} else {
		// Create new collaboration
		collaboration := models.Collaboration{
			NoteID:     uint(noteID),
			UserID:     targetUser.ID,
			Permission: req.Permission,
		}
		config.DB.Create(&collaboration)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note shared successfully"})
}

func GetCollaborators(c *gin.Context) {
	userID := c.GetUint("user_id")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	// Check if user has access to the note
	var note models.Note
	query := config.DB.Where("id = ?", noteID).
		Where("owner_id = ? OR id IN (SELECT note_id FROM collaborations WHERE user_id = ? AND deleted_at IS NULL)", userID, userID)
	
	if err := query.First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	var collaborations []models.Collaboration
	if err := config.DB.Where("note_id = ?", noteID).Preload("User").Find(&collaborations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch collaborators"})
		return
	}

	c.JSON(http.StatusOK, collaborations)
}
