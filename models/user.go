package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
	Phone    string `json:"phone" gorm:"unique"`
	
	StudentID  string `json:"student_id" gorm:"unique;not null"`
	Course     string `json:"course"`
	Department string `json:"department"`
	Year       int    `json:"year"`
	
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Role      string  `json:"role" gorm:"default:'user'"`
	
	IsOwner        bool    `json:"is_owner" gorm:"default:false"`
	OwnerRating    float64 `json:"owner_rating" gorm:"default:0"`
	TotalVehicles  int     `json:"total_vehicles" gorm:"default:0"`
	
	RenterRating   float64 `json:"renter_rating" gorm:"default:0"`
	TotalRentals   int     `json:"total_rentals" gorm:"default:0"`
	
	DrivingLicense string    `json:"driving_license"`
	LicenseNumber  string    `json:"license_number"`
	LicenseExpiry  time.Time `json:"license_expiry"`
	AadharCard     string    `json:"aadhar_card"`
	
	IsVerified    bool       `json:"is_verified" gorm:"default:false"`
	VerifiedAt    *time.Time `json:"verified_at"`
	
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	LastActive    time.Time `json:"last_active"`
	
	UpiID         string `json:"upi_id"`
	BankAccount   string `json:"bank_account"`
}
