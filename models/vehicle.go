package models

import (
	"time"
	"gorm.io/gorm"
)

type Vehicle struct {
	gorm.Model
	
	OwnerID       uint   `json:"owner_id"`
	Owner         User   `json:"owner" gorm:"foreignKey:OwnerID"`
	
	VehicleType   string `json:"vehicle_type"`
	Brand         string `json:"brand"`
	Model         string `json:"model"`
	Year          int    `json:"year"`
	Color         string `json:"color"`
	VehicleNumber string `json:"vehicle_number" gorm:"unique;not null"`
	
	PricePerHour  float64 `json:"price_per_hour"`
	PricePerDay   float64 `json:"price_per_day"`
	MinRentalHours int    `json:"min_rental_hours" gorm:"default:1"`
	MaxRentalDays  int    `json:"max_rental_days" gorm:"default:7"`
	
	HasHelmet     bool   `json:"has_helmet"`
	FuelType      string `json:"fuel_type"`
	Transmission  string `json:"transmission"`
	Mileage       float64 `json:"mileage"`
	SeatingCapacity int  `json:"seating_capacity"`
	
	RCDocument    string    `json:"rc_document"`
	Insurance     string    `json:"insurance"`
	InsuranceExpiry time.Time `json:"insurance_expiry"`
	IsVerified    bool      `json:"is_verified" gorm:"default:false"`
	
	Location      string  `json:"location"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	
	IsAvailable   bool    `json:"is_available" gorm:"default:true"`
	IsActive      bool    `json:"is_active" gorm:"default:true"`
	
	Rating        float64 `json:"rating" gorm:"default:0"`
	TotalBookings int     `json:"total_bookings" gorm:"default:0"`
	
	Images        string  `json:"images" gorm:"type:text"`
	Description   string  `json:"description" gorm:"type:text"`
	Rules         string  `json:"rules" gorm:"type:text"`
}
