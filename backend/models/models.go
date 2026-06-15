package models

import "time"

const (
	SeatAvailable = "AVAILABLE"
	SeatLocked    = "LOCKED"
	SeatBooked    = "BOOKED"

	BookingConfirmed = "CONFIRMED"
	BookingCancelled = "CANCELLED"
	BookingTimeout   = "TIMEOUT"

	RoleUser  = "USER"
	RoleAdmin = "ADMIN"
)

type Seat struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	ShowtimeID string `json:"showtime_id"`
	LockedBy   string `json:"locked_by,omitempty"`
}

type Showtime struct {
	ID         string `json:"id" bson:"_id"`
	Movie      string `json:"movie" bson:"movie"`
	Date       string `json:"date" bson:"date"`
	Time       string `json:"time" bson:"time"`
	TotalSeats int    `json:"total_seats" bson:"total_seats"`
}

type Booking struct {
	ID         string    `json:"id" bson:"_id"`
	UserID     string    `json:"user_id" bson:"user_id"`
	UserEmail  string    `json:"user_email,omitempty" bson:"user_email,omitempty"`
	SeatID     string    `json:"seat_id" bson:"seat_id"`
	ShowtimeID string    `json:"showtime_id" bson:"showtime_id"`
	Movie      string    `json:"movie" bson:"movie"`
	Date       string    `json:"date" bson:"date"`
	Status     string    `json:"status" bson:"status"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
}

type AuditLog struct {
	ID         string    `json:"id" bson:"_id"`
	Event      string    `json:"event" bson:"event"`
	UserID     string    `json:"user_id,omitempty" bson:"user_id,omitempty"`
	SeatID     string    `json:"seat_id,omitempty" bson:"seat_id,omitempty"`
	ShowtimeID string    `json:"showtime_id,omitempty" bson:"showtime_id,omitempty"`
	Message    string    `json:"message" bson:"message"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
}

type BookingEvent struct {
	Type       string `json:"type"`
	BookingID  string `json:"booking_id"`
	UserID     string `json:"user_id"`
	UserEmail  string `json:"user_email"`
	SeatID     string `json:"seat_id"`
	ShowtimeID string `json:"showtime_id"`
	Movie      string `json:"movie"`
}

type WSMessage struct {
	Type       string `json:"type"`
	SeatID     string `json:"seat_id"`
	ShowtimeID string `json:"showtime_id"`
	Status     string `json:"status"`
	LockedBy   string `json:"locked_by,omitempty"`
}