package routes

import (
	"proj/handlers"
	"proj/middleware"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	api := r.Group("/api")

	api.POST("/register", handlers.Register)
	api.POST("/login", handlers.Login)

	protected := api.Group("")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/profile", handlers.GetProfile)
		protected.GET("/users", handlers.GetUsers)

		protected.POST("/vehicles", handlers.CreateVehicle)
		protected.GET("/my-vehicles", handlers.GetMyVehicles)
		protected.PUT("/vehicles/:id", handlers.UpdateVehicle)
		protected.DELETE("/vehicles/:id", handlers.DeleteVehicle)

		protected.POST("/vehicles/:id/availability", handlers.SetAvailability)
		protected.PUT("/availability/:id", handlers.UpdateAvailability)
		protected.DELETE("/availability/:id", handlers.DeleteAvailability)

		protected.POST("/bookings", handlers.CreateBooking)
		protected.GET("/bookings", handlers.GetBookings)
		protected.GET("/bookings/active", handlers.GetActiveBooking)
		protected.GET("/bookings/history", handlers.GetBookingHistory)
		protected.GET("/bookings/:id", handlers.GetBookingByID)
		protected.POST("/bookings/:id/confirm", handlers.ConfirmBooking)
		protected.POST("/bookings/:id/cancel", handlers.CancelBooking)
		protected.POST("/bookings/:id/pickup/generate-otp", handlers.GeneratePickupOTP)
		protected.POST("/bookings/:id/pickup/verify-otp", handlers.VerifyPickupOTP)
		protected.POST("/bookings/:id/return/generate-otp", handlers.GenerateReturnOTP)
		protected.POST("/bookings/:id/return/verify-otp", handlers.VerifyReturnOTP)

		protected.POST("/documents", handlers.UploadDocument)
		protected.GET("/documents", handlers.GetMyDocuments)
		protected.GET("/documents/:id", handlers.GetDocumentByID)
		protected.DELETE("/documents/:id", handlers.DeleteDocument)

		protected.GET("/admin/documents/pending", handlers.GetPendingDocuments)
		protected.POST("/admin/documents/:id/verify", handlers.VerifyDocument)
	}

	api.GET("/vehicles", handlers.GetVehicles)
	api.GET("/vehicles/:id", handlers.GetVehicleByID)
	api.GET("/vehicles/:id/availability", handlers.GetAvailability)
	api.GET("/availability/check", handlers.CheckAvailability)
}