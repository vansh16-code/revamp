package models

import (
	"time"
	"gorm.io/gorm"
)

type Availability struct {
	gorm.Model
	VehicleID     uint      `json:"vehicle_id"`
	Vehicle       Vehicle   `json:"vehicle" gorm:"foreignKey:VehicleID"`
	
	AvailableFrom time.Time `json:"available_from"`
	AvailableTo   time.Time `json:"available_to"`
	
	IsRecurring   bool   `json:"is_recurring" gorm:"default:false"`
	DaysOfWeek    string `json:"days_of_week"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	
	Status        string `json:"status" gorm:"default:'available'"`
}
