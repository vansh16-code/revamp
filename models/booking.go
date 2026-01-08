package models

import (
	"time"
	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	
	VehicleID     uint    `json:"vehicle_id"`
	Vehicle       Vehicle `json:"vehicle" gorm:"foreignKey:VehicleID"`
	RenterID      uint    `json:"renter_id"`
	Renter        User    `json:"renter" gorm:"foreignKey:RenterID"`
	OwnerID       uint    `json:"owner_id"`
	Owner         User    `json:"owner" gorm:"foreignKey:OwnerID"`
	
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Duration      int       `json:"duration"`
	
	PricePerHour  float64 `json:"price_per_hour"`
	TotalPrice    float64 `json:"total_price"`
	SecurityDeposit float64 `json:"security_deposit"`
	
	PickupLocation  string    `json:"pickup_location"`
	ReturnLocation  string    `json:"return_location"`
	PickupTime      time.Time `json:"pickup_time"`
	ReturnTime      time.Time `json:"return_time"`
	
	Status        string `json:"status" gorm:"default:'pending'"`
	
	PickupOTP     string `json:"pickup_otp"`
	ReturnOTP     string `json:"return_otp"`
	
	FuelLevelStart   int    `json:"fuel_level_start"`
	FuelLevelEnd     int    `json:"fuel_level_end"`
	OdometerStart    int    `json:"odometer_start"`
	OdometerEnd      int    `json:"odometer_end"`
	DamageReportStart string `json:"damage_report_start" gorm:"type:text"`
	DamageReportEnd   string `json:"damage_report_end" gorm:"type:text"`
	
	PickupImages  string `json:"pickup_images" gorm:"type:text"`
	ReturnImages  string `json:"return_images" gorm:"type:text"`
	
	Notes         string `json:"notes" gorm:"type:text"`
}
