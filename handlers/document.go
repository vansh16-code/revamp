package handlers

import (
	"net/http"
	"time"

	"proj/config"
	"proj/models"

	"github.com/gin-gonic/gin"
)

type UploadDocumentRequest struct {
	DocumentType   string `json:"document_type" binding:"required"`
	DocumentNumber string `json:"document_number"`
	DocumentURL    string `json:"document_url" binding:"required"`
	IssueDate      string `json:"issue_date"`
	ExpiryDate     string `json:"expiry_date"`
	VehicleID      *uint  `json:"vehicle_id"`
	Notes          string `json:"notes"`
}

func UploadDocument(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in context"})
		return
	}

	var req UploadDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validTypes := map[string]bool{
		models.DocumentTypeDrivingLicense: true,
		models.DocumentTypeStudentID:      true,
		models.DocumentTypeAadhar:         true,
		models.DocumentTypeRC:             true,
		models.DocumentTypeInsurance:      true,
		models.DocumentTypePUC:            true,
	}

	if !validTypes[req.DocumentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document type"})
		return
	}

	if req.VehicleID != nil {
		var vehicle models.Vehicle
		if err := config.DB.First(&vehicle, *req.VehicleID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}

		if vehicle.OwnerID != uid {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't own this vehicle"})
			return
		}
	}

	document := models.Document{
		UserID:         &uid,
		VehicleID:      req.VehicleID,
		DocumentType:   req.DocumentType,
		DocumentNumber: req.DocumentNumber,
		DocumentURL:    req.DocumentURL,
		Status:         models.DocumentStatusPending,
		Notes:          req.Notes,
	}

	if req.IssueDate != "" {
		issueDate, err := time.Parse("2006-01-02", req.IssueDate)
		if err == nil {
			document.IssueDate = &issueDate
		}
	}

	if req.ExpiryDate != "" {
		expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
		if err == nil {
			document.ExpiryDate = &expiryDate
		}
	}

	if err := config.DB.Create(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload document"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Document uploaded successfully",
		"document": document,
	})
}

func GetMyDocuments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in context"})
		return
	}

	documentType := c.Query("type")
	status := c.Query("status")

	query := config.DB.Model(&models.Document{}).Where("user_id = ?", uid)

	if documentType != "" {
		query = query.Where("document_type = ?", documentType)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var documents []models.Document
	if err := query.Order("created_at DESC").Find(&documents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     len(documents),
		"documents": documents,
	})
}

func GetDocumentByID(c *gin.Context) {
	documentID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in context"})
		return
	}

	var document models.Document
	if err := config.DB.Preload("User").Preload("Vehicle").First(&document, documentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	if document.UserID != nil && *document.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this document"})
		return
	}

	c.JSON(http.StatusOK, document)
}

func DeleteDocument(c *gin.Context) {
	documentID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in context"})
		return
	}

	var document models.Document
	if err := config.DB.First(&document, documentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	if document.UserID != nil && *document.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this document"})
		return
	}

	if document.Status == models.DocumentStatusApproved {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete approved document"})
		return
	}

	if err := config.DB.Delete(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document deleted successfully",
	})
}

func GetPendingDocuments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	documentType := c.Query("type")

	query := config.DB.Model(&models.Document{}).Where("status = ?", models.DocumentStatusPending)

	if documentType != "" {
		query = query.Where("document_type = ?", documentType)
	}

	var documents []models.Document
	if err := query.Preload("User").Preload("Vehicle").Order("created_at ASC").Find(&documents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     len(documents),
		"documents": documents,
	})
}

type VerifyDocumentRequest struct {
	Status          string `json:"status" binding:"required"`
	RejectionReason string `json:"rejection_reason"`
}

func VerifyDocument(c *gin.Context) {
	documentID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in context"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var req VerifyDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status != models.DocumentStatusApproved && req.Status != models.DocumentStatusRejected {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be 'approved' or 'rejected'"})
		return
	}

	if req.Status == models.DocumentStatusRejected && req.RejectionReason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rejection reason is required"})
		return
	}

	var document models.Document
	if err := config.DB.First(&document, documentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	if document.Status != models.DocumentStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document is not pending verification"})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":           req.Status,
		"verified_by":      uid,
		"verified_at":      &now,
		"rejection_reason": req.RejectionReason,
	}

	if err := config.DB.Model(&document).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify document"})
		return
	}

	if req.Status == models.DocumentStatusApproved {
		updateVerificationStatus(&document)
	}

	if err := config.DB.Preload("User").Preload("Vehicle").First(&document, documentID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Document verified but failed to load details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Document verified successfully",
		"document": document,
	})
}

func updateVerificationStatus(document *models.Document) {
	if document.UserID != nil {
		switch document.DocumentType {
		case models.DocumentTypeDrivingLicense:
			config.DB.Model(&models.User{}).Where("id = ?", *document.UserID).Update("license_verified", true)
		case models.DocumentTypeAadhar:
			config.DB.Model(&models.User{}).Where("id = ?", *document.UserID).Update("aadhar_verified", true)
		case models.DocumentTypeStudentID:
			config.DB.Model(&models.User{}).Where("id = ?", *document.UserID).Update("student_id_verified", true)
		}

		var user models.User
		config.DB.First(&user, *document.UserID)
		if user.LicenseVerified && user.AadharVerified && user.StudentIDVerified {
			now := time.Now()
			config.DB.Model(&user).Updates(map[string]interface{}{
				"is_verified": true,
				"verified_at": &now,
			})
		}
	}

	if document.VehicleID != nil {
		switch document.DocumentType {
		case models.DocumentTypeRC:
			config.DB.Model(&models.Vehicle{}).Where("id = ?", *document.VehicleID).Update("rc_verified", true)
		case models.DocumentTypeInsurance:
			config.DB.Model(&models.Vehicle{}).Where("id = ?", *document.VehicleID).Update("insurance_verified", true)
		case models.DocumentTypePUC:
			config.DB.Model(&models.Vehicle{}).Where("id = ?", *document.VehicleID).Update("puc_verified", true)
		}

		var vehicle models.Vehicle
		config.DB.First(&vehicle, *document.VehicleID)
		if vehicle.RCVerified && vehicle.InsuranceVerified {
			now := time.Now()
			config.DB.Model(&vehicle).Updates(map[string]interface{}{
				"is_verified": true,
				"verified_at": &now,
			})
		}
	}
}
