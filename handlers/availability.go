package handlers

import (
	"net/http"
	"time"

	"proj/config"
	"proj/models"

	"github.com/gin-gonic/gin"
)

type SetAvailabilityRequest struct {
	AvailableFrom time.Time `json:"available_from" binding:"required"`
	AvailableTo   time.Time `json:"available_to" binding:"required"`
	IsRecurring   bool      `json:"is_recurring"`
	DaysOfWeek    string    `json:"days_of_week"`
	StartTime     string    `json:"start_time"`
	EndTime       string    `json:"end_time"`
}

type UpdateAvailabilityRequest struct {
	AvailableFrom *time.Time `json:"available_from"`
	AvailableTo   *time.Time `json:"available_to"`
	Status        *string    `json:"status"`
}

func SetAvailability(c *gin.Context) {
	vehicleID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var vehicle models.Vehicle
	if err := config.DB.First(&vehicle, vehicleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}

	if vehicle.OwnerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't own this vehicle"})
		return
	}

	var req SetAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.AvailableFrom.After(req.AvailableTo) || req.AvailableFrom.Equal(req.AvailableTo) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "available_from must be before available_to"})
		return
	}

	if req.AvailableFrom.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot set availability in the past"})
		return
	}

	var overlappingBookings int64
	result := config.DB.Model(&models.Booking{}).
		Where("vehicle_id = ? AND status IN ? AND ((start_time BETWEEN ? AND ?) OR (end_time BETWEEN ? AND ?) OR (start_time <= ? AND end_time >= ?))",
			vehicleID,
			[]string{models.BookingStatusConfirmed, models.BookingStatusOngoing},
			req.AvailableFrom, req.AvailableTo,
			req.AvailableFrom, req.AvailableTo,
			req.AvailableFrom, req.AvailableTo,
		).Count(&overlappingBookings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for overlapping bookings"})
		return
	}

	if overlappingBookings > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Vehicle has confirmed bookings during this time"})
		return
	}

	availability := models.Availability{
		VehicleID:     vehicle.ID,
		AvailableFrom: req.AvailableFrom,
		AvailableTo:   req.AvailableTo,
		IsRecurring:   req.IsRecurring,
		DaysOfWeek:    req.DaysOfWeek,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Status:        models.AvailabilityStatusAvailable,
	}

	if err := config.DB.Create(&availability).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set availability"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Availability set successfully",
		"availability": availability,
	})
}

func GetAvailability(c *gin.Context) {
	vehicleID := c.Param("id")

	var vehicle models.Vehicle
	if err := config.DB.First(&vehicle, vehicleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}

	status := c.Query("status")
	fromDate := c.Query("from")
	toDate := c.Query("to")

	query := config.DB.Where("vehicle_id = ?", vehicleID)

	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status = ?", models.AvailabilityStatusAvailable)
	}

	if fromDate != "" && toDate != "" {
		from, err1 := time.Parse(time.RFC3339, fromDate)
		to, err2 := time.Parse(time.RFC3339, toDate)
		if err1 == nil && err2 == nil {
			query = query.Where("available_from <= ? AND available_to >= ?", to, from)
		}
	}

	query = query.Where("available_to >= ?", time.Now())

	var availabilities []models.Availability
	if err := query.Order("available_from ASC").Find(&availabilities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch availability"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vehicle_id":     vehicle.ID,
		"total":          len(availabilities),
		"availabilities": availabilities,
	})
}

func UpdateAvailability(c *gin.Context) {
	availabilityID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var availability models.Availability
	if err := config.DB.Preload("Vehicle").First(&availability, availabilityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Availability not found"})
		return
	}

	if availability.Vehicle.OwnerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't own this vehicle"})
		return
	}

	var req UpdateAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})

	if req.AvailableFrom != nil {
		if req.AvailableFrom.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot set availability in the past"})
			return
		}
		updates["available_from"] = *req.AvailableFrom
	}

	if req.AvailableTo != nil {
		if req.AvailableTo.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot set availability in the past"})
			return
		}
		updates["available_to"] = *req.AvailableTo
	}

	if req.AvailableFrom != nil && req.AvailableTo != nil {
		if !req.AvailableFrom.Before(*req.AvailableTo) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "available_from must be before available_to"})
			return
		}

		var conflictingBookings int64
		result := config.DB.Model(&models.Booking{}).
			Where("vehicle_id = ? AND status IN ? AND start_time <= ? AND end_time >= ?",
				availability.VehicleID,
				[]string{models.BookingStatusConfirmed, models.BookingStatusOngoing},
				*req.AvailableTo,
				*req.AvailableFrom,
			).Count(&conflictingBookings)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for booking conflicts"})
			return
		}

		if conflictingBookings > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Cannot update: vehicle has confirmed bookings during this time"})
			return
		}
	}

	if req.Status != nil {
		validStatuses := []string{
			models.AvailabilityStatusAvailable,
			models.AvailabilityStatusBlocked,
			models.AvailabilityStatusMaintenance,
		}
		isValid := false
		for _, status := range validStatuses {
			if *req.Status == status {
				isValid = true
				break
			}
		}
		if !isValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	if err := config.DB.Model(&availability).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update availability"})
		return
	}

	config.DB.Preload("Vehicle").First(&availability, availabilityID)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Availability updated successfully",
		"availability": availability,
	})
}

func DeleteAvailability(c *gin.Context) {
	availabilityID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var availability models.Availability
	if err := config.DB.Preload("Vehicle").First(&availability, availabilityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Availability not found"})
		return
	}

	if availability.Vehicle.OwnerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't own this vehicle"})
		return
	}

	var activeBookings int64
	result := config.DB.Model(&models.Booking{}).
		Where("vehicle_id = ? AND status IN ? AND start_time <= ? AND end_time >= ?",
			availability.VehicleID,
			[]string{models.BookingStatusConfirmed, models.BookingStatusOngoing},
			availability.AvailableTo,
			availability.AvailableFrom,
		).Count(&activeBookings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for active bookings"})
		return
	}

	if activeBookings > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete availability with active bookings"})
		return
	}

	if err := config.DB.Delete(&availability).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete availability"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Availability deleted successfully",
	})
}

func CheckAvailability(c *gin.Context) {
	vehicleID := c.Query("vehicle_id")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	if vehicleID == "" || startTimeStr == "" || endTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_id, start_time, and end_time are required"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format. Use RFC3339 (e.g., 2026-01-10T09:00:00Z)"})
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format. Use RFC3339 (e.g., 2026-01-10T18:00:00Z)"})
		return
	}

	if startTime.After(endTime) || startTime.Equal(endTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_time must be before end_time"})
		return
	}

	if startTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot check availability in the past"})
		return
	}

	var vehicle models.Vehicle
	if err := config.DB.First(&vehicle, vehicleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}

	if !vehicle.IsAvailable || !vehicle.IsActive {
		c.JSON(http.StatusOK, gin.H{
			"available": false,
			"reason":    "Vehicle is not available",
		})
		return
	}

	var availabilitySlots int64
	result := config.DB.Model(&models.Availability{}).
		Where("vehicle_id = ? AND status = ? AND available_from <= ? AND available_to >= ?",
			vehicleID,
			models.AvailabilityStatusAvailable,
			startTime,
			endTime,
		).Count(&availabilitySlots)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check availability slots"})
		return
	}

	if availabilitySlots == 0 {
		c.JSON(http.StatusOK, gin.H{
			"available": false,
			"reason":    "No availability slot for this time range",
		})
		return
	}

	var conflictingBookings int64
	result = config.DB.Model(&models.Booking{}).
		Where("vehicle_id = ? AND status IN ? AND start_time <= ? AND end_time >= ?",
			vehicleID,
			[]string{models.BookingStatusConfirmed, models.BookingStatusOngoing},
			endTime,
			startTime,
		).Count(&conflictingBookings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for booking conflicts"})
		return
	}

	if conflictingBookings > 0 {
		c.JSON(http.StatusOK, gin.H{
			"available": false,
			"reason":    "Vehicle is already booked for this time",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"available":   true,
		"vehicle_id":  vehicle.ID,
		"start_time":  startTime,
		"end_time":    endTime,
		"price_info": gin.H{
			"price_per_km":   vehicle.PricePerKm,
			"price_per_hour": vehicle.PricePerHour,
			"price_per_day":  vehicle.PricePerDay,
			"base_price":     vehicle.BasePrice,
		},
	})
}
