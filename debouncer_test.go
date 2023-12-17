package debouncer

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"os"
	"testing"
	"time"
)

var client *redis.Client

func Init() {
	client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
}
func TestDebounce(t *testing.T) {
	Init()
	key := "debounce:test"
	thread := 30
	for i := 0; i < thread; i++ {
		go func(ii int) {
			rand.Seed(time.Now().UnixNano())
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			Debounce(client, key, time.Second*2, time.Second*5, func() {
				fmt.Println("callback ", ii)
			})
		}(i)
	}
	time.Sleep(time.Second * 10)
}
