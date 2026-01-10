package handlers

import (
	"net/http"
	"strconv"

	"proj/config"
	"proj/models"

	"github.com/gin-gonic/gin"
)

type CreateVehicleRequest struct {
	VehicleType     string  `json:"vehicle_type" binding:"required"`
	Brand           string  `json:"brand" binding:"required"`
	VehicleModel    string  `json:"vehicle_model" binding:"required"`
	Year            int     `json:"year" binding:"required"`
	Color           string  `json:"color"`
	VehicleNumber   string  `json:"vehicle_number" binding:"required"`
	PricePerKm      int64   `json:"price_per_km"`
	PricePerHour    int64   `json:"price_per_hour"`
	PricePerDay     int64   `json:"price_per_day"`
	BasePrice       int64   `json:"base_price"`
	MinRentalHours  int     `json:"min_rental_hours"`
	MaxRentalDays   int     `json:"max_rental_days"`
	HasHelmet       bool    `json:"has_helmet"`
	FuelType        string  `json:"fuel_type"`
	Transmission    string  `json:"transmission"`
	Mileage         float64 `json:"mileage"`
	SeatingCapacity int     `json:"seating_capacity"`
	Location        string  `json:"location" binding:"required"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Description     string  `json:"description"`
	Rules           string  `json:"rules"`
}

type UpdateVehicleRequest struct {
	PricePerKm      *int64   `json:"price_per_km"`
	PricePerHour    *int64   `json:"price_per_hour"`
	PricePerDay     *int64   `json:"price_per_day"`
	BasePrice       *int64   `json:"base_price"`
	MinRentalHours  *int     `json:"min_rental_hours"`
	MaxRentalDays   *int     `json:"max_rental_days"`
	HasHelmet       *bool    `json:"has_helmet"`
	Location        *string  `json:"location"`
	Latitude        *float64 `json:"latitude"`
	Longitude       *float64 `json:"longitude"`
	Description     *string  `json:"description"`
	Rules           *string  `json:"rules"`
	IsAvailable     *bool    `json:"is_available"`
}

func CreateVehicle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.DrivingLicense == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please upload driving license first"})
		return
	}

	var existingVehicle models.Vehicle
	if err := config.DB.Where("vehicle_number = ?", req.VehicleNumber).First(&existingVehicle).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Vehicle with this number already registered"})
		return
	}

	vehicle := models.Vehicle{
		OwnerID:         userID.(uint),
		VehicleType:     req.VehicleType,
		Brand:           req.Brand,
		VehicleModel:    req.VehicleModel,
		Year:            req.Year,
		Color:           req.Color,
		VehicleNumber:   req.VehicleNumber,
		PricePerKm:      req.PricePerKm,
		PricePerHour:    req.PricePerHour,
		PricePerDay:     req.PricePerDay,
		BasePrice:       req.BasePrice,
		MinRentalHours:  req.MinRentalHours,
		HasHelmet:       req.HasHelmet,
		FuelType:        req.FuelType,
		Transmission:    req.Transmission,
		Mileage:         req.Mileage,
		SeatingCapacity: req.SeatingCapacity,
		Location:        req.Location,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		Description:     req.Description,
		Rules:           req.Rules,
		IsAvailable:     true,
		IsActive:        true,
	}

	if err := config.DB.Create(&vehicle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vehicle"})
		return
	}

	if err := config.DB.Model(&user).Updates(map[string]interface{}{
		"is_owner":       true,
		"total_vehicles": user.TotalVehicles + 1,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user stats"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Vehicle listed successfully",
		"vehicle": vehicle,
	})
}

func GetVehicles(c *gin.Context) {
	vehicleType := c.Query("type")
	maxPrice := c.Query("max_price")
	minRating := c.Query("min_rating")
	location := c.Query("location")

	query := config.DB.Where("is_active = ? AND is_available = ?", true, true)

	if vehicleType != "" {
		query = query.Where("vehicle_type = ?", vehicleType)
	}

	if maxPrice != "" {
		price, err := strconv.ParseInt(maxPrice, 10, 64)
		if err == nil {
			query = query.Where("price_per_day <= ?", price)
		}
	}

	if minRating != "" {
		rating, err := strconv.ParseFloat(minRating, 64)
		if err == nil {
			query = query.Where("rating >= ?", rating)
		}
	}

	if location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	var vehicles []models.Vehicle
	if err := query.Preload("Owner").Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    len(vehicles),
		"vehicles": vehicles,
	})
}

func GetVehicleByID(c *gin.Context) {
	vehicleID := c.Param("id")

	var vehicle models.Vehicle
	if err := config.DB.Preload("Owner").First(&vehicle, vehicleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}

	var availability []models.Availability
	config.DB.Where("vehicle_id = ? AND status = ?", vehicleID, models.AvailabilityStatusAvailable).
		Find(&availability)

	c.JSON(http.StatusOK, gin.H{
		"vehicle":      vehicle,
		"availability": availability,
	})
}

func UpdateVehicle(c *gin.Context) {
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

	var req UpdateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})

	if req.PricePerKm != nil {
		updates["price_per_km"] = *req.PricePerKm
	}
	if req.PricePerHour != nil {
		updates["price_per_hour"] = *req.PricePerHour
	}
	if req.PricePerDay != nil {
		updates["price_per_day"] = *req.PricePerDay
	}
	if req.BasePrice != nil {
		updates["base_price"] = *req.BasePrice
	}
	if req.MinRentalHours != nil {
		updates["min_rental_hours"] = *req.MinRentalHours
	}
	if req.MaxRentalDays != nil {
		updates["max_rental_days"] = *req.MaxRentalDays
	}
	if req.HasHelmet != nil {
		updates["has_helmet"] = *req.HasHelmet
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.Latitude != nil {
		updates["latitude"] = *req.Latitude
	}
	if req.Longitude != nil {
		updates["longitude"] = *req.Longitude
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Rules != nil {
		updates["rules"] = *req.Rules
	}
	if req.IsAvailable != nil {
		updates["is_available"] = *req.IsAvailable
	}

	if err := config.DB.Model(&vehicle).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vehicle updated successfully",
		"vehicle": vehicle,
	})
}

func DeleteVehicle(c *gin.Context) {
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

	var activeBookings int64
	config.DB.Model(&models.Booking{}).
		Where("vehicle_id = ? AND status IN ?", vehicleID, []string{
			models.BookingStatusConfirmed,
			models.BookingStatusOngoing,
		}).
		Count(&activeBookings)

	if activeBookings > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete vehicle with active bookings"})
		return
	}

	if err := config.DB.Model(&vehicle).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vehicle"})
		return
	}

	var user models.User
	config.DB.First(&user, userID)
	if user.TotalVehicles > 0 {
		config.DB.Model(&user).Update("total_vehicles", user.TotalVehicles-1)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vehicle deleted successfully",
	})
}

func GetMyVehicles(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var vehicles []models.Vehicle
	if err := config.DB.Where("owner_id = ? AND is_active = ?", userID, true).
		Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    len(vehicles),
		"vehicles": vehicles,
	})
}
