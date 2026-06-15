package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cinema/config"

	"github.com/redis/go-redis/v9"
)

const LockTTL = 5 * time.Minute

var (
	ctx = context.Background()
	Rdb *redis.Client
)

func InitRedis(cfg *config.Config) error {
	Rdb = redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	if err := Rdb.Ping(ctx).Err(); err != nil {
		return err
	}

	// Enable keyspace notifications for lock expiry (best-effort).
	_ = Rdb.ConfigSet(ctx, "notify-keyspace-events", "Ex").Err()
	return nil
}

func lockKey(showtimeID, seatID string) string {
	return fmt.Sprintf("lock:%s:%s", showtimeID, seatID)
}

func LockSeat(showtimeID, seatID, userID string) (bool, error) {
	ok, err := Rdb.SetNX(ctx, lockKey(showtimeID, seatID), userID, LockTTL).Result()
	return ok, err
}

func UnlockSeat(showtimeID, seatID string) error {
	return Rdb.Del(ctx, lockKey(showtimeID, seatID)).Err()
}

func GetLockOwner(showtimeID, seatID string) (string, error) {
	val, err := Rdb.Get(ctx, lockKey(showtimeID, seatID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func IsSeatLocked(showtimeID, seatID string) bool {
	owner, _ := GetLockOwner(showtimeID, seatID)
	return owner != ""
}

func SubscribeLockExpiry(handler func(showtimeID, seatID string)) {
	pubsub := Rdb.PSubscribe(ctx, "__keyevent@*__:expired")
	go func() {
		for msg := range pubsub.Channel() {
			showtimeID, seatID := parseLockKey(msg.Payload)
			if showtimeID != "" && seatID != "" {
				handler(showtimeID, seatID)
			}
		}
	}()
}

func parseLockKey(key string) (string, string) {
	const prefix = "lock:"
	if !strings.HasPrefix(key, prefix) {
		return "", ""
	}
	rest := key[len(prefix):]
	idx := strings.Index(rest, ":")
	if idx < 0 {
		return "", ""
	}
	return rest[:idx], rest[idx+1:]
}
