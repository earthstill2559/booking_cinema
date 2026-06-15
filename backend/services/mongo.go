package services

import (
	"context"
	"time"

	"cinema/config"
	"cinema/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient  *mongo.Client
	bookingsCol  *mongo.Collection
	auditLogsCol *mongo.Collection
	showtimesCol *mongo.Collection
	defaultMovie = "Avatar 3"
	defaultDate  = "2026-06-12"
	defaultTime  = "19:00"
)

func InitMongo(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	mongoClient = client
	db := client.Database("cinema")
	bookingsCol = db.Collection("bookings")
	auditLogsCol = db.Collection("audit_logs")
	showtimesCol = db.Collection("showtimes")

	_, _ = bookingsCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "showtime_id", Value: 1}, {Key: "seat_id", Value: 1}},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{"status": models.BookingConfirmed}),
	})
	_, _ = auditLogsCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: -1}},
	})

	return nil
}

// DefaultShowtime returns the showtime configured as cfg.DefaultShowtimeID,
// falling back to the hardcoded defaults if it hasn't been created in Mongo yet.
func DefaultShowtime(cfg *config.Config) models.Showtime {
	showtime, err := GetShowtimeByID(cfg, cfg.DefaultShowtimeID)
	if err != nil {
		return models.Showtime{
			ID:         cfg.DefaultShowtimeID,
			Movie:      defaultMovie,
			Date:       defaultDate,
			Time:       defaultTime,
			TotalSeats: len(AllSeatIDs()),
		}
	}
	return showtime
}

// GetShowtimeByID fetches a showtime from Mongo. If the id matches
// cfg.DefaultShowtimeID and no document exists yet, it returns the
// hardcoded default (for backward compatibility before any showtime
// has been created via the admin panel).
func GetShowtimeByID(cfg *config.Config, id string) (models.Showtime, error) {
	ctx := context.Background()

	var showtime models.Showtime
	err := showtimesCol.FindOne(ctx, bson.M{"_id": id}).Decode(&showtime)
	if err == mongo.ErrNoDocuments {
		if id == cfg.DefaultShowtimeID {
			return models.Showtime{
				ID:         cfg.DefaultShowtimeID,
				Movie:      defaultMovie,
				Date:       defaultDate,
				Time:       defaultTime,
				TotalSeats: len(AllSeatIDs()),
			}, nil
		}
		return models.Showtime{}, err
	}
	if err != nil {
		return models.Showtime{}, err
	}
	return showtime, nil
}

func CreateShowtime(showtime *models.Showtime) error {
	ctx := context.Background()
	_, err := showtimesCol.InsertOne(ctx, showtime)
	return err
}

func ListShowtimes() ([]models.Showtime, error) {
	ctx := context.Background()
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
	cursor, err := showtimesCol.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.Showtime
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
func ListShowtimesForUsers(cfg *config.Config) ([]models.Showtime, error) {
	showtimes, err := ListShowtimes()
	if err != nil {
		return nil, err
	}
	if len(showtimes) == 0 {
		return []models.Showtime{DefaultShowtime(cfg)}, nil
	}
	return showtimes, nil
}

func IsSeatBooked(showtimeID, seatID string) (bool, error) {
	ctx := context.Background()
	count, err := bookingsCol.CountDocuments(ctx, bson.M{
		"showtime_id": showtimeID,
		"seat_id":     seatID,
		"status":      models.BookingConfirmed,
	})
	return count > 0, err
}

func CreateBooking(booking *models.Booking) error {
	ctx := context.Background()
	booking.CreatedAt = time.Now()
	_, err := bookingsCol.InsertOne(ctx, booking)
	return err
}

func ListBookings(filter bson.M) ([]models.Booking, error) {
	ctx := context.Background()
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := bookingsCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.Booking
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func SaveAuditLog(log *models.AuditLog) error {
	ctx := context.Background()
	log.CreatedAt = time.Now()
	_, err := auditLogsCol.InsertOne(ctx, log)
	return err
}

func ListAuditLogs(filter bson.M) ([]models.AuditLog, error) {
	ctx := context.Background()
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(100)
	cursor, err := auditLogsCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.AuditLog
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}