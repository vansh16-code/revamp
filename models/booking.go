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
	DurationHours int       `json:"duration_hours"`
	
	PricingModel  string `json:"pricing_model" gorm:"default:'distance'"` 
	
	PricePerKm      int64 `json:"price_per_km"`
	PricePerHour    int64 `json:"price_per_hour"`
	BasePrice       int64 `json:"base_price"`
	EstimatedPrice  int64 `json:"estimated_price"`
	FinalPrice      int64 `json:"final_price"`
	SecurityDeposit int64 `json:"security_deposit"`
	
	PickupLocation  string    `json:"pickup_location"`
	ReturnLocation  string    `json:"return_location"`
	PickupTime      time.Time `json:"pickup_time"`
	ReturnTime      time.Time `json:"return_time"`
	
	Status        string `json:"status" gorm:"default:'pending'"`
	
	PickupOTP     string `json:"pickup_otp"`
	ReturnOTP     string `json:"return_otp"`
	
	OdometerStartKm       int     `json:"odometer_start_km"`
	OdometerEndKm         int     `json:"odometer_end_km"`
	ActualDistanceKm      float64 `json:"actual_distance_km"`
	
	FuelLevelStartPercent int     `json:"fuel_level_start_percent"`
	FuelLevelEndPercent   int     `json:"fuel_level_end_percent"`
	FuelConsumedLiters    float64 `json:"fuel_consumed_liters"`
	FuelCostCharged       int64   `json:"fuel_cost_charged"`
	
	DamageReportStart     string `json:"damage_report_start" gorm:"type:text"`
	DamageReportEnd       string `json:"damage_report_end" gorm:"type:text"`
	
	PickupImages  string `json:"pickup_images" gorm:"type:text"`
	ReturnImages  string `json:"return_images" gorm:"type:text"`
	
	HasOBDData    bool   `json:"has_obd_data" gorm:"default:false"`
	OBDTrackerID  *uint  `json:"obd_tracker_id"`
	
	Notes         string `json:"notes" gorm:"type:text"`
}
