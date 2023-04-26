package limiter

import (
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func TestName(t *testing.T) {

}

func TestExampleNewLimiter(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.161.129:6379",
		Password: "UOCk68SDQTOXNPsW",
		DB:       3,
	})
	_ = rdb.FlushDB().Err()

	limiter := NewLimiter(rdb)

	for {
		res, err := limiter.Allow("project:123", PerSecond(10))
		if err != nil {
			panic(err)
		}
		t.Log("allowed", res.Allowed, "remaining", res.Remaining)
		if res.Allowed != 1 {
			time.Sleep(10 * time.Millisecond)
		}
		time.Sleep(70 * time.Millisecond)
	}
}
