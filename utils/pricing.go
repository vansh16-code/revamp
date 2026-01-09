package utils

import (
	"math"
	"proj/models"
	"time"
)

func CalculateDistanceBasedPrice(distanceKm float64, pricePerKm int64, basePrice int64) int64 {
	return basePrice + int64(math.Round(distanceKm*float64(pricePerKm)))
}

func CalculateTimeBasedPrice(durationHours int, pricePerHour int64, basePrice int64) int64 {
	return basePrice + (int64(durationHours) * pricePerHour)
}

func CalculateHybridPrice(distanceKm float64, durationHours int, pricePerKm int64, pricePerHour int64, basePrice int64) int64 {
	distancePrice := int64(math.Round(distanceKm * float64(pricePerKm)))
	timePrice := int64(durationHours) * pricePerHour
	return basePrice + distancePrice + timePrice
}

func CalculateFuelCost(fuelConsumedLiters float64, fuelPricePerLiter int64) int64 {
	return int64(math.Round(fuelConsumedLiters * float64(fuelPricePerLiter)))
}

func EstimatePrice(vehicle *models.Vehicle, estimatedDistanceKm float64, durationHours int, pricingModel string) int64 {
	switch pricingModel {
	case models.PricingModelDistance:
		return CalculateDistanceBasedPrice(estimatedDistanceKm, vehicle.PricePerKm, vehicle.BasePrice)
	case models.PricingModelTime:
		return CalculateTimeBasedPrice(durationHours, vehicle.PricePerHour, vehicle.BasePrice)
	case models.PricingModelHybrid:
		return CalculateHybridPrice(estimatedDistanceKm, durationHours, vehicle.PricePerKm, vehicle.PricePerHour, vehicle.BasePrice)
	default:
		return CalculateDistanceBasedPrice(estimatedDistanceKm, vehicle.PricePerKm, vehicle.BasePrice)
	}
}

func CalculateFinalPrice(booking *models.Booking) int64 {
	var finalPrice int64
	
	switch booking.PricingModel {
	case models.PricingModelDistance:
		finalPrice = CalculateDistanceBasedPrice(booking.ActualDistanceKm, booking.PricePerKm, booking.BasePrice)
	case models.PricingModelTime:
		actualHours := int(math.Ceil(booking.ReturnTime.Sub(booking.PickupTime).Hours()))
		finalPrice = CalculateTimeBasedPrice(actualHours, booking.PricePerHour, booking.BasePrice)
	case models.PricingModelHybrid:
		actualHours := int(math.Ceil(booking.ReturnTime.Sub(booking.PickupTime).Hours()))
		finalPrice = CalculateHybridPrice(booking.ActualDistanceKm, actualHours, booking.PricePerKm, booking.PricePerHour, booking.BasePrice)
	default:
		finalPrice = booking.EstimatedPrice
	}
	
	finalPrice += booking.FuelCostCharged
	
	return finalPrice
}

func CalculateDuration(startTime, endTime time.Time) int {
	return int(math.Ceil(endTime.Sub(startTime).Hours()))
}
