package handlers

import (
	"math"
	"net/http"
	"time"

	"proj/config"
	"proj/models"
	"proj/utils"

	"github.com/gin-gonic/gin"
)

type CreateBookingRequest struct {
	VehicleID       uint      `json:"vehicle_id" binding:"required"`
	StartTime       time.Time `json:"start_time" binding:"required"`
	EndTime         time.Time `json:"end_time" binding:"required"`
	PickupLocation  string    `json:"pickup_location" binding:"required"`
	ReturnLocation  string    `json:"return_location"`
	PricingModel    string    `json:"pricing_model"`
	EstimatedDistanceKm float64 `json:"estimated_distance_km"`
	Notes           string    `json:"notes"`
}

func CreateBooking(c *gin.Context) {
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

	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.StartTime.After(req.EndTime) || req.StartTime.Equal(req.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_time must be before end_time"})
		return
	}

	if req.StartTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot book in the past"})
		return
	}

	var vehicle models.Vehicle
	if err := config.DB.Preload("Owner").First(&vehicle, req.VehicleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}

	if !vehicle.IsAvailable || !vehicle.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vehicle is not available"})
		return
	}

	if vehicle.OwnerID == uid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot book your own vehicle"})
		return
	}

	var availabilitySlots int64
	result := config.DB.Model(&models.Availability{}).
		Where("vehicle_id = ? AND status = ? AND available_from <= ? AND available_to >= ?",
			req.VehicleID,
			models.AvailabilityStatusAvailable,
			req.StartTime,
			req.EndTime,
		).Count(&availabilitySlots)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check availability"})
		return
	}

	if availabilitySlots == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Vehicle is not available for this time range"})
		return
	}

	var conflictingBookings int64
	result = config.DB.Model(&models.Booking{}).
		Where("vehicle_id = ? AND status IN ? AND start_time <= ? AND end_time >= ?",
			req.VehicleID,
			[]string{models.BookingStatusConfirmed, models.BookingStatusOngoing},
			req.EndTime,
			req.StartTime,
		).Count(&conflictingBookings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for conflicts"})
		return
	}

	if conflictingBookings > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Vehicle is already booked for this time"})
		return
	}

	durationHours := int(math.Ceil(req.EndTime.Sub(req.StartTime).Hours()))

	pricingModel := req.PricingModel
	if pricingModel == "" {
		pricingModel = models.PricingModelDistance
	}

	estimatedPrice := utils.EstimatePrice(&vehicle, req.EstimatedDistanceKm, durationHours, pricingModel)

	securityDeposit := vehicle.PricePerDay

	returnLocation := req.ReturnLocation
	if returnLocation == "" {
		returnLocation = req.PickupLocation
	}

	booking := models.Booking{
		VehicleID:       req.VehicleID,
		RenterID:        uid,
		OwnerID:         vehicle.OwnerID,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		DurationHours:   durationHours,
		PricingModel:    pricingModel,
		PricePerKm:      vehicle.PricePerKm,
		PricePerHour:    vehicle.PricePerHour,
		BasePrice:       vehicle.BasePrice,
		EstimatedPrice:  estimatedPrice,
		SecurityDeposit: securityDeposit,
		PickupLocation:  req.PickupLocation,
		ReturnLocation:  returnLocation,
		Status:          models.BookingStatusPending,
		Notes:           req.Notes,
	}

	if err := config.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	result := config.DB.Preload("Vehicle").Preload("Owner").Preload("Renter").First(&booking, booking.ID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking created but failed to load details"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Booking created successfully",
		"booking": booking,
	})
}

func GetBookings(c *gin.Context) {
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

	role := c.Query("role")
	status := c.Query("status")

	query := config.DB.Model(&models.Booking{})

	if role == "owner" {
		query = query.Where("owner_id = ?", uid)
	} else if role == "renter" {
		query = query.Where("renter_id = ?", uid)
	} else {
		query = query.Where("renter_id = ? OR owner_id = ?", uid, uid)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var bookings []models.Booking
	if err := query.Preload("Vehicle").Preload("Owner").Preload("Renter").
		Order("created_at DESC").Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    len(bookings),
		"bookings": bookings,
	})
}

func GetBookingByID(c *gin.Context) {
	bookingID := c.Param("id")
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

	var booking models.Booking
	if err := config.DB.Preload("Vehicle").Preload("Owner").Preload("Renter").
		First(&booking, bookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if booking.RenterID != uid && booking.OwnerID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this booking"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func ConfirmBooking(c *gin.Context) {
	bookingID := c.Param("id")
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

	var booking models.Booking
	if err := config.DB.First(&booking, bookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if booking.OwnerID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the vehicle owner can confirm bookings"})
		return
	}

	if booking.Status != models.BookingStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking is not in pending status"})
		return
	}

	if err := config.DB.Model(&booking).Update("status", models.BookingStatusConfirmed).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm booking"})
		return
	}

	result := config.DB.Preload("Vehicle").Preload("Owner").Preload("Renter").First(&booking, bookingID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking confirmed but failed to load details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking confirmed successfully",
		"booking": booking,
	})
}

func CancelBooking(c *gin.Context) {
	bookingID := c.Param("id")
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

	var booking models.Booking
	if err := config.DB.First(&booking, bookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if booking.RenterID != uid && booking.OwnerID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to cancel this booking"})
		return
	}

	if booking.Status == models.BookingStatusOngoing {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot cancel ongoing booking"})
		return
	}

	if booking.Status == models.BookingStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot cancel completed booking"})
		return
	}

	if booking.Status == models.BookingStatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking is already cancelled"})
		return
	}

	if err := config.DB.Model(&booking).Update("status", models.BookingStatusCancelled).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel booking"})
		return
	}

	result := config.DB.Preload("Vehicle").Preload("Owner").Preload("Renter").First(&booking, bookingID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking cancelled but failed to load details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking cancelled successfully",
		"booking": booking,
	})
}

func GetActiveBooking(c *gin.Context) {
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

	var booking models.Booking
	err := config.DB.Preload("Vehicle").Preload("Owner").Preload("Renter").
		Where("renter_id = ? AND status = ?", uid, models.BookingStatusOngoing).
		First(&booking).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"active_booking": nil,
			"message":        "No active booking",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"active_booking": booking,
	})
}

func GetBookingHistory(c *gin.Context) {
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

	var bookings []models.Booking
	if err := config.DB.Preload("Vehicle").Preload("Owner").Preload("Renter").
		Where("renter_id = ? AND status IN ?", uid, []string{
			models.BookingStatusCompleted,
			models.BookingStatusCancelled,
		}).
		Order("created_at DESC").
		Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch booking history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    len(bookings),
		"bookings": bookings,
	})
}
