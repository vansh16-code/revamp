package models

import (
	"time"
	"gorm.io/gorm"
)

type Vehicle struct {
	gorm.Model
	
	OwnerID       uint   `json:"owner_id" gorm:"not null"`
	Owner         User   `json:"owner" gorm:"foreignKey:OwnerID"`
	
	VehicleType   string `json:"vehicle_type"`
	Brand         string `json:"brand"`
	VehicleModel  string `json:"vehicle_model"`
	Year          int    `json:"year"`
	Color         string `json:"color"`
	VehicleNumber string `json:"vehicle_number" gorm:"unique;not null"`
	
	PricePerKm    int64 `json:"price_per_km"`
	PricePerHour  int64 `json:"price_per_hour"`
	PricePerDay   int64 `json:"price_per_day"`
	BasePrice     int64 `json:"base_price"`
	MinRentalHours int  `json:"min_rental_hours" gorm:"default:1"`
	MaxRentalDays  int  `json:"max_rental_days" gorm:"default:7"`
	
	HasHelmet     bool    `json:"has_helmet"`
	FuelType      string  `json:"fuel_type"`
	Transmission  string  `json:"transmission"`
	Mileage       float64 `json:"mileage"`
	SeatingCapacity int   `json:"seating_capacity"`
	
	RCDocument      string     `json:"rc_document"`
	RCNumber        string     `json:"rc_number"`
	RCVerified      bool       `json:"rc_verified" gorm:"default:false"`
	
	Insurance       string     `json:"insurance"`
	InsuranceNumber string     `json:"insurance_number"`
	InsuranceExpiry *time.Time `json:"insurance_expiry,omitempty"`
	InsuranceVerified bool     `json:"insurance_verified" gorm:"default:false"`
	
	PUCDocument     string     `json:"puc_document"`
	PUCExpiry       *time.Time `json:"puc_expiry,omitempty"`
	PUCVerified     bool       `json:"puc_verified" gorm:"default:false"`
	
	IsVerified      bool       `json:"is_verified" gorm:"default:false"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`
	
	Location      string  `json:"location"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	
	HasOBDTracker bool   `json:"has_obd_tracker" gorm:"default:false"`
	OBDTrackerID  *uint  `json:"obd_tracker_id"`
	
	IsAvailable   bool    `json:"is_available" gorm:"default:true"`
	IsActive      bool    `json:"is_active" gorm:"default:true"`
	
	Rating        float64 `json:"rating" gorm:"default:0"`
	TotalBookings int     `json:"total_bookings" gorm:"default:0"`
	TotalKmDriven int     `json:"total_km_driven" gorm:"default:0"`
	
	Images        string  `json:"images" gorm:"type:text"`
	Description   string  `json:"description" gorm:"type:text"`
	Rules         string  `json:"rules" gorm:"type:text"`
}
