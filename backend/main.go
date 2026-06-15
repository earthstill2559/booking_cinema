package main

import (
	"log"
	"net/http"

	"cinema/config"
	"cinema/handlers"
	"cinema/middleware"
	"cinema/models"
	"cinema/services"
	"cinema/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if err := services.InitRedis(cfg); err != nil {
		log.Fatalf("redis init failed: %v", err)
	}
	if err := services.InitMongo(cfg); err != nil {
		log.Fatalf("mongo init failed: %v", err)
	}
	if err := services.InitMQ(cfg); err != nil {
		log.Fatalf("rabbitmq init failed: %v", err)
	}

	services.StartBookingConsumer(services.HandleBookingEventAsync)

	services.SubscribeLockExpiry(func(showtimeID, seatID string) {
		services.LogBookingTimeout(showtimeID, seatID)
		services.LogSeatReleased("", seatID, showtimeID, "Lock expired automatically")
		ws.Broadcast(showtimeID, seatID, models.SeatAvailable, "")
	})

	bookingHandler := &handlers.BookingHandler{Cfg: cfg}
	adminHandler := &handlers.AdminHandler{Cfg: cfg}

	r := gin.Default()
	r.Use(corsMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/ws", ws.HandleWebSocket)

		api := r.Group("/api")
	{
		api.GET("/showtime", bookingHandler.GetShowtime)
		api.GET("/showtimes", bookingHandler.ListShowtimes)
		api.GET("/seats", bookingHandler.GetSeats)

		auth := api.Group("", middleware.AuthRequired(cfg))
		{
			auth.GET("/me", bookingHandler.Me)
			auth.POST("/seats/:seat_id/lock", bookingHandler.LockSeat)
			auth.POST("/seats/:seat_id/unlock", bookingHandler.UnlockSeat)
			auth.POST("/bookings/confirm", bookingHandler.ConfirmBooking)
		}

		admin := api.Group("/admin", middleware.AuthRequired(cfg), middleware.AdminRequired())
		{
			admin.GET("/bookings", adminHandler.ListBookings)
			admin.GET("/audit-logs", adminHandler.ListAuditLogs)
			admin.POST("/showtimes", adminHandler.CreateShowtime)
			admin.GET("/showtimes", adminHandler.ListShowtimes)
		}
	}

	// Legacy routes for backward compatibility
	r.GET("/seats", bookingHandler.GetSeats)

	log.Printf("Server starting on :%s (dev_mode=%v)", cfg.Port, cfg.DevMode)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
