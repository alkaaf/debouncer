package debouncer

import (
	"github.com/go-redis/redis/v8"
	"time"
)

func Debounce(c *redis.Client, key string, duration time.Duration, atomicLock time.Duration, callback func()) {
	ctx := c.Context()
	ownTicket := int64(0)
	err := c.Watch(ctx, func(tx *redis.Tx) error {
		var err error
		// get the current ticket
		ownTicket, err = tx.Incr(ctx, key).Result()
		if err != nil {
			return err
		}
		_, err = tx.Expire(ctx, key, duration+(time.Second*5)).Result()
		if err != nil {
			return err
		}
		return nil
	}, key)
	if err != nil {
		return
	}
	// run the the debounced function in the background
	go func() {
		time.Sleep(duration)
		var err error
		currentTicket, err := c.Get(ctx, key).Int64()
		if err != nil {
			return
		}

		// if current ticket is greater than own ticket, it means that the key has been updated by another thread
		if currentTicket > ownTicket {
			return
		}
		lock := key + "_lock"
		b, err := c.SetNX(ctx, lock, 1, atomicLock).Result()
		if err != nil {
			return
		}
		if !b {
			return
		}
		callback()
		c.Del(ctx, lock)

	}()
}
