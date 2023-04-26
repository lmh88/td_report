package redis

import (
	rediS "github.com/go-redis/redis"
	"math/rand"
	"sync/atomic"
	"time"

	red "github.com/go-redis/redis"
	"td_report/common/stringx"
	"td_report/pkg/logger"
)

const (
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
	randomLen = 16
)

// A RedisLock is a redis lock. 参考https://github.com/zeromicro/go-zero/blob/master/core/stores/redis/redislock.go
type RedisLock struct {
	store   *rediS.Client
	seconds uint32
	count   int32
	key     string
	id      string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewRedisLock returns a RedisLock.
func NewRedisLock(store *rediS.Client, key string) *RedisLock {
	return &RedisLock{
		store: store,
		key:   key,
		id:    stringx.Randn(randomLen),
	}
}

// Acquire acquires the lock.
func (rl *RedisLock) Acquire() (bool, error) {
	newCount := atomic.AddInt32(&rl.count, 1)
	if newCount > 1 {
		return true, nil
	}

	seconds := atomic.LoadUint32(&rl.seconds)
	time.Sleep(1 * time.Second)
	//ok, err := rl.store.SetnxEx(rl.key, rl.id, int(seconds+1)) // +1s for tolerance
	ok, err := rl.store.SetNX(rl.key, rl.id, time.Duration(seconds+1)*time.Second).Result() // +1s for tolerance
	if err == red.Nil {
		atomic.AddInt32(&rl.count, -1)
		return false, nil
	} else if err != nil {
		atomic.AddInt32(&rl.count, -1)
		logger.Errorf("Error on acquiring lock for %s, %s", rl.key, err.Error())
		return false, err
	} else if !ok {
		atomic.AddInt32(&rl.count, -1)
		return false, nil
	}

	return true, nil
}

// Release releases the lock.
func (rl *RedisLock) Release() (bool, error) {
	newCount := atomic.AddInt32(&rl.count, -1)
	if newCount > 0 {
		return true, nil
	}

	//resp, err := rl.store.Eval(delCommand, []string{rl.key}, []string{rl.id})
	resp, err := rl.store.Eval(delCommand, []string{rl.key}, []string{rl.id}).Result()
	if err != nil {
		return false, err
	}

	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}

	return reply == 1, nil
}

// SetExpire sets the expiration.
func (rl *RedisLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.seconds, uint32(seconds))
}
