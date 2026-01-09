package models

const (
	AvailabilityStatusAvailable    = "available"
	AvailabilityStatusBooked       = "booked"
	AvailabilityStatusBlocked      = "blocked"
	AvailabilityStatusMaintenance  = "maintenance"
)

const (
	BookingStatusPending   = "pending"
	BookingStatusConfirmed = "confirmed"
	BookingStatusOngoing   = "ongoing"
	BookingStatusCompleted = "completed"
	BookingStatusCancelled = "cancelled"
	BookingStatusDisputed  = "disputed"
)

const (
	PricingModelDistance = "distance"
	PricingModelTime     = "time"
	PricingModelHybrid   = "hybrid"
)
