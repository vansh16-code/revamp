package models

import (
	"time"
	"gorm.io/gorm"
)

type OBDTracker struct {
	gorm.Model
	VehicleID     uint    `json:"vehicle_id" gorm:"not null"`
	Vehicle       Vehicle `json:"vehicle" gorm:"foreignKey:VehicleID"`
	
	DeviceID      string `json:"device_id" gorm:"unique;not null"`
	DeviceModel   string `json:"device_model"`
	IsActive      bool   `json:"is_active" gorm:"default:true"`
	LastSyncedAt  *time.Time `json:"last_synced_at,omitempty"`
	
	BatteryLevel  int    `json:"battery_level"`
	SignalStrength int   `json:"signal_strength"`
}

type OBDReading struct {
	gorm.Model
	BookingID     uint    `json:"booking_id" gorm:"not null"`
	Booking       Booking `json:"booking" gorm:"foreignKey:BookingID"`
	OBDTrackerID  uint    `json:"obd_tracker_id" gorm:"not null"`
	OBDTracker    OBDTracker `json:"obd_tracker" gorm:"foreignKey:OBDTrackerID"`
	
	Timestamp     time.Time `json:"timestamp" gorm:"not null"`
	
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Speed         float64 `json:"speed"`
	
	OdometerKm    int     `json:"odometer_km"`
	TripDistanceKm float64 `json:"trip_distance_km"`
	
	FuelLevelPercent int     `json:"fuel_level_percent"`
	FuelConsumedLiters float64 `json:"fuel_consumed_liters"`
	
	EngineRPM     int     `json:"engine_rpm"`
	EngineTemp    int     `json:"engine_temp"`
	CoolantTemp   int     `json:"coolant_temp"`
	
	ThrottlePosition int  `json:"throttle_position"`
	IsEngineOn    bool    `json:"is_engine_on"`
	IsMoving      bool    `json:"is_moving"`
	
	BatteryVoltage float64 `json:"battery_voltage"`
	
	DiagnosticCodes string `json:"diagnostic_codes" gorm:"type:text"`
}
