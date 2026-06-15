package services

import (
	"strconv"

	"cinema/config"
	"cinema/models"
)

const seatsPerRow = 8

var seatRows = []string{"A", "B", "C", "D", "E"}
var seatCols = []int{1, 2, 3, 4, 5, 6, 7, 8}

// AllSeatIDs returns the original fixed 40-seat layout (A1-E8),
// kept as the fallback for the legacy default showtime.
func AllSeatIDs() []string {
	ids := make([]string, 0, len(seatRows)*len(seatCols))
	for _, row := range seatRows {
		for _, col := range seatCols {
			ids = append(ids, row+strconv.Itoa(col))
		}
	}
	return ids
}

// SeatIDsForCount generates seat IDs (e.g. A1..A8, B1..B8, ...)
// for the given total seat count, 8 seats per row.
func SeatIDsForCount(total int) []string {
	if total <= 0 {
		return AllSeatIDs()
	}

	ids := make([]string, 0, total)
	rows := (total + seatsPerRow - 1) / seatsPerRow
	count := 0
	for r := 0; r < rows && count < total; r++ {
		rowLabel := string(rune('A' + r))
		for col := 1; col <= seatsPerRow && count < total; col++ {
			ids = append(ids, rowLabel+strconv.Itoa(col))
			count++
		}
	}
	return ids
}

func GetSeatsForShowtime(cfg *config.Config, showtimeID string) ([]models.Seat, error) {
	showtime, err := GetShowtimeByID(cfg, showtimeID)
	if err != nil {
		return nil, err
	}

	ids := SeatIDsForCount(showtime.TotalSeats)
	seats := make([]models.Seat, 0, len(ids))

	for _, id := range ids {
		seat := models.Seat{
			ID:         id,
			ShowtimeID: showtimeID,
			Status:     models.SeatAvailable,
		}

		booked, err := IsSeatBooked(showtimeID, id)
		if err != nil {
			return nil, err
		}
		if booked {
			seat.Status = models.SeatBooked
			seats = append(seats, seat)
			continue
		}

		owner, err := GetLockOwner(showtimeID, id)
		if err != nil {
			return nil, err
		}
		if owner != "" {
			seat.Status = models.SeatLocked
			seat.LockedBy = owner
		}

		seats = append(seats, seat)
	}

	return seats, nil
}