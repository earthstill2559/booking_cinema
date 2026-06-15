package handlers

import (
	"fmt"
	"net/http"

	"cinema/config"
	"cinema/models"
	"cinema/services"
	"cinema/ws"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingHandler struct {
	Cfg *config.Config
}

type lockRequest struct {
	ShowtimeID string `json:"showtime_id"`
}

type confirmRequest struct {
	ShowtimeID string   `json:"showtime_id"`
	SeatIDs    []string `json:"seat_ids"`
	Movie      string   `json:"movie,omitempty"`
	Date       string   `json:"date,omitempty"`
}

func (h *BookingHandler) GetShowtime(c *gin.Context) {
	c.JSON(http.StatusOK, services.DefaultShowtime(h.Cfg))
}

func (h *BookingHandler) GetSeats(c *gin.Context) {
	showtimeID := c.DefaultQuery("showtime_id", h.Cfg.DefaultShowtimeID)
	seats, err := services.GetSeatsForShowtime(h.Cfg, showtimeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, seats)
}

func (h *BookingHandler) LockSeat(c *gin.Context) {
	seatID := c.Param("seat_id")
	userID := c.GetString("user_id")
	showtimeID := h.Cfg.DefaultShowtimeID

	var body lockRequest
	_ = c.ShouldBindJSON(&body)
	if body.ShowtimeID != "" {
		showtimeID = body.ShowtimeID
	}

	booked, err := services.IsSeatBooked(showtimeID, seatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if booked {
		c.JSON(http.StatusConflict, gin.H{"error": "seat already booked"})
		return
	}

	owner, err := services.GetLockOwner(showtimeID, seatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if owner != "" && owner != userID {
		c.JSON(http.StatusConflict, gin.H{"error": "seat locked by another user"})
		return
	}
	if owner == userID {
		c.JSON(http.StatusOK, gin.H{"message": "already locked", "seat_id": seatID, "status": models.SeatLocked})
		return
	}

	ok, err := services.LockSeat(showtimeID, seatID, userID)
	if err != nil {
		services.LogLockFail(userID, seatID, showtimeID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lock failed"})
		return
	}
	if !ok {
		services.LogLockFail(userID, seatID, showtimeID, fmt.Errorf("seat already locked"))
		c.JSON(http.StatusConflict, gin.H{"error": "seat locked by another user"})
		return
	}

	ws.Broadcast(showtimeID, seatID, models.SeatLocked, userID)
	c.JSON(http.StatusOK, gin.H{"message": "locked", "seat_id": seatID, "status": models.SeatLocked})
}

func (h *BookingHandler) UnlockSeat(c *gin.Context) {
	seatID := c.Param("seat_id")
	userID := c.GetString("user_id")
	showtimeID := h.Cfg.DefaultShowtimeID

	var body lockRequest
	_ = c.ShouldBindJSON(&body)
	if body.ShowtimeID != "" {
		showtimeID = body.ShowtimeID
	}

	owner, err := services.GetLockOwner(showtimeID, seatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if owner == "" {
		c.JSON(http.StatusOK, gin.H{"message": "seat not locked"})
		return
	}
	if owner != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not lock owner"})
		return
	}

	if err := services.UnlockSeat(showtimeID, seatID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	services.LogSeatReleased(userID, seatID, showtimeID, "User cancelled selection")
	ws.Broadcast(showtimeID, seatID, models.SeatAvailable, "")
	c.JSON(http.StatusOK, gin.H{"message": "released", "seat_id": seatID, "status": models.SeatAvailable})
}

func (h *BookingHandler) ConfirmBooking(c *gin.Context) {
	userID := c.GetString("user_id")
	userEmail := c.GetString("user_email")
	showtimeID := h.Cfg.DefaultShowtimeID

	var body confirmRequest
	if err := c.ShouldBindJSON(&body); err != nil || len(body.SeatIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "seat_ids required"})
		return
	}
	if body.ShowtimeID != "" {
		showtimeID = body.ShowtimeID
	}

	showtime := services.DefaultShowtime(h.Cfg)
	if body.Movie != "" {
		showtime.Movie = body.Movie
	}
	if body.Date != "" {
		showtime.Date = body.Date
	}
	confirmed := make([]string, 0, len(body.SeatIDs))

	for _, seatID := range body.SeatIDs {
		owner, err := services.GetLockOwner(showtimeID, seatID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if owner != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "seat " + seatID + " not locked by you"})
			return
		}

		booked, err := services.IsSeatBooked(showtimeID, seatID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if booked {
			c.JSON(http.StatusConflict, gin.H{"error": "seat " + seatID + " already booked"})
			return
		}

		booking := &models.Booking{
			ID:         uuid.New().String(),
			UserID:     userID,
			UserEmail:  userEmail,
			SeatID:     seatID,
			ShowtimeID: showtimeID,
			Movie:      showtime.Movie,
			Date:       showtime.Date,
			Status:     models.BookingConfirmed,
		}
		if err := services.CreateBooking(booking); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "double booking prevented for seat " + seatID})
			return
		}

		_ = services.UnlockSeat(showtimeID, seatID)
		services.LogBookingSuccess(userID, seatID, showtimeID)

		event := models.BookingEvent{
			Type:       "booking.success",
			BookingID:  booking.ID,
			UserID:     userID,
			UserEmail:  userEmail,
			SeatID:     seatID,
			ShowtimeID: showtimeID,
			Movie:      showtime.Movie,
		}
		_ = services.PublishBookingSuccess(event)

		ws.Broadcast(showtimeID, seatID, models.SeatBooked, "")
		confirmed = append(confirmed, seatID)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "booking confirmed",
		"seat_ids": confirmed,
		"status":   models.SeatBooked,
	})
}

func (h *BookingHandler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user_id": c.GetString("user_id"),
		"email":   c.GetString("user_email"),
		"role":    c.GetString("role"),
	})
}

func (h *BookingHandler) ListShowtimes(c *gin.Context) {
	showtimes, err := services.ListShowtimesForUsers(h.Cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, showtimes)
}