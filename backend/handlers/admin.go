package handlers

import (
	"net/http"

	"cinema/config"
	"cinema/models"
	"cinema/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type AdminHandler struct {
	Cfg *config.Config
}

func (h *AdminHandler) ListBookings(c *gin.Context) {
	filter := bson.M{"status": "CONFIRMED"}

	if movie := c.Query("movie"); movie != "" {
		filter["movie"] = movie
	}
	if date := c.Query("date"); date != "" {
		filter["date"] = date
	}
	if userID := c.Query("user_id"); userID != "" {
		filter["user_id"] = userID
	}

	bookings, err := services.ListBookings(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *AdminHandler) ListAuditLogs(c *gin.Context) {
	filter := bson.M{}
	if event := c.Query("event"); event != "" {
		filter["event"] = event
	}

	logs, err := services.ListAuditLogs(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

type createShowtimeRequest struct {
	ID         string `json:"id,omitempty"`
	Movie      string `json:"movie" binding:"required"`
	Date       string `json:"date" binding:"required"`
	Time       string `json:"time" binding:"required"`
	TotalSeats int    `json:"total_seats" binding:"required,min=1"`
}

func (h *AdminHandler) CreateShowtime(c *gin.Context) {
	var body createShowtimeRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := body.ID
	if id == "" {
		id = uuid.New().String()
	}

	showtime := &models.Showtime{
		ID:         id,
		Movie:      body.Movie,
		Date:       body.Date,
		Time:       body.Time,
		TotalSeats: body.TotalSeats,
	}

	if err := services.CreateShowtime(showtime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, showtime)
}

func (h *AdminHandler) ListShowtimes(c *gin.Context) {
	showtimes, err := services.ListShowtimes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, showtimes)
}