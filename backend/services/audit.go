package services

import (
	"fmt"
	"log"

	"cinema/models"

	"github.com/google/uuid"
)

const (
	EventBookingSuccess = "BOOKING_SUCCESS"
	EventBookingTimeout = "BOOKING_TIMEOUT"
	EventSeatReleased   = "SEAT_RELEASED"
	EventSystemError    = "SYSTEM_ERROR"
)

func LogAudit(event, userID, seatID, showtimeID, message string) {
	entry := &models.AuditLog{
		ID:         uuid.New().String(),
		Event:      event,
		UserID:     userID,
		SeatID:     seatID,
		ShowtimeID: showtimeID,
		Message:    message,
	}
	if err := SaveAuditLog(entry); err != nil {
		log.Printf("audit log failed: %v", err)
	}
}

func LogLockFail(userID, seatID, showtimeID string, err error) {
	LogAudit(EventSystemError, userID, seatID, showtimeID, fmt.Sprintf("Lock fail: %v", err))
}

func LogBookingSuccess(userID, seatID, showtimeID string) {
	LogAudit(EventBookingSuccess, userID, seatID, showtimeID, "Booking confirmed successfully")
}

func LogBookingTimeout(showtimeID, seatID string) {
	LogAudit(EventBookingTimeout, "", seatID, showtimeID, "Lock expired after 5 minutes without payment")
}

func LogSeatReleased(userID, seatID, showtimeID, reason string) {
	LogAudit(EventSeatReleased, userID, seatID, showtimeID, reason)
}

func HandleBookingEventAsync(event models.BookingEvent) {
	LogAudit(EventBookingSuccess, event.UserID, event.SeatID, event.ShowtimeID,
		fmt.Sprintf("Async notification for booking %s (mock email to %s)", event.BookingID, event.UserEmail))
	log.Printf("[NOTIFICATION MOCK] Booking success: user=%s seat=%s movie=%s",
		event.UserID, event.SeatID, event.Movie)
}
