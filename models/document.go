package models

import (
	"time"
	"gorm.io/gorm"
)

const (
	DocumentTypeDrivingLicense = "driving_license"
	DocumentTypeStudentID      = "student_id"
	DocumentTypeAadhar         = "aadhar"
	DocumentTypeRC             = "rc"
	DocumentTypeInsurance      = "insurance"
	DocumentTypePUC            = "puc"
	
	DocumentStatusPending   = "pending"
	DocumentStatusApproved  = "approved"
	DocumentStatusRejected  = "rejected"
	DocumentStatusExpired   = "expired"
)

type Document struct {
	gorm.Model
	
	UserID      *uint  `json:"user_id"`
	User        *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	
	VehicleID   *uint    `json:"vehicle_id"`
	Vehicle     *Vehicle `json:"vehicle,omitempty" gorm:"foreignKey:VehicleID"`
	
	DocumentType   string `json:"document_type" gorm:"not null"`
	DocumentNumber string `json:"document_number"`
	DocumentURL    string `json:"document_url" gorm:"not null"`
	
	IssueDate      *time.Time `json:"issue_date,omitempty"`
	ExpiryDate     *time.Time `json:"expiry_date,omitempty"`
	
	Status         string     `json:"status" gorm:"default:'pending'"`
	VerifiedBy     *uint      `json:"verified_by"`
	VerifiedAt     *time.Time `json:"verified_at,omitempty"`
	RejectionReason string    `json:"rejection_reason" gorm:"type:text"`
	
	Notes          string     `json:"notes" gorm:"type:text"`
}
