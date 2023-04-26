package redis

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"td_report/common/stringx"
)

func TestRedisLock(t *testing.T) {
	Redisclient := NewDefaultRds()
	client := Redisclient.GetClient()
	key := stringx.Rand()
	firstLock := NewRedisLock(client, key)
	firstLock.SetExpire(5)
	firstAcquire, err := firstLock.Acquire()
	assert.Nil(t, err)
	assert.True(t, firstAcquire)

	secondLock := NewRedisLock(client, key)
	secondLock.SetExpire(5)
	againAcquire, err := secondLock.Acquire()
	assert.Nil(t, err)
	assert.False(t, againAcquire)

	release, err := firstLock.Release()
	assert.Nil(t, err)
	assert.True(t, release)

	endAcquire, err := secondLock.Acquire()
	assert.Nil(t, err)
	assert.True(t, endAcquire)

	endAcquire, err = secondLock.Acquire()
	assert.Nil(t, err)
	assert.True(t, endAcquire)

	release, err = secondLock.Release()
	assert.Nil(t, err)
	assert.True(t, release)

	againAcquire, err = firstLock.Acquire()
	assert.Nil(t, err)
	assert.False(t, againAcquire)

	release, err = secondLock.Release()
	assert.Nil(t, err)
	assert.True(t, release)

	firstAcquire, err = firstLock.Acquire()
	assert.Nil(t, err)
	assert.True(t, firstAcquire)
}

func TestRedisLock1(t *testing.T) {
	wg := sync.WaitGroup{}
	n := 0
	Redisclient := NewDefaultRds()
	client := Redisclient.GetClient()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			key := stringx.Rand()
			firstLock := NewRedisLock(client, key)
			firstLock.SetExpire(5)
			_, err := firstLock.Acquire()
			if err != nil {
				t.Log("error ")
			}
			n = n + num
			firstLock.Release()
		}(i)
	}

	wg.Wait()
	fmt.Println(n)
}
