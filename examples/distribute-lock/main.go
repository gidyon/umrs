package main

import (
	"context"
	"fmt"
	"github.com/gidyon/micros/pkg/conn"
	"github.com/go-redis/redis"
	"math/rand"
	"sync"
	"time"
)

const (
	tokenSet         = "tokens"
	tokenLeaseTime   = 10 * time.Second
	maximumLeaseTime = 30 * time.Second
)

var (
	// errLeaseTimeExceeded ...
	errLeaseTimeExceeded = fmt.Errorf("maximum lease time for token is %d seconds", 10)
)

type lockerPayload struct {
	lockID  string
	client  *redis.Client
	channel <-chan *redis.Message
}

func acquire(ctx context.Context, locker *lockerPayload) error {
	if t, ok := ctx.Deadline(); ok {
		rem := t.Sub(time.Now())
		if rem > tokenLeaseTime {
			return errLeaseTimeExceeded
		}
	}
	ch := make(chan struct{})
	var err error
	go func() {
		err = acquireLock(ctx, locker)
		close(ch)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return err
	}
}

func acquireLock(ctx context.Context, locker *lockerPayload) error {
	v, err := locker.client.SAdd(tokenSet, locker.lockID).Result()
	if err != nil {
		return fmt.Errorf("failed to acquire token: %w", err)
	}

	if v == 1 {
		return nil
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to acquire token: %v", err)
	case <-locker.channel:
		return acquireLock(ctx, locker)
	}
}

type Locker interface {
	Acquire(context.Context) error
	Release(context.Context) error
	Released() <-chan struct{}
}

type redisLocker struct {
	client  *redis.Client
	lockID  string
	channel <-chan *redis.Message
}

// NewRedisLocker creates a singleton distributed locker
func NewRedisLocker(client *redis.Client) Locker {
	return &redisLocker{}
}

func (rl *redisLocker) acquireLock(ctx context.Context) error {
	v, err := rl.client.SAdd(tokenSet, rl.lockID).Result()
	if err != nil {
		return fmt.Errorf("failed to acquire token: %w", err)
	}

	if v == 1 {
		return nil
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to acquire token: %v", err)
	case <-rl.channel:
		return rl.acquireLock(ctx)
	}
}

func (rl *redisLocker) Acquire(ctx context.Context) error {
	if t, ok := ctx.Deadline(); ok {
		rem := t.Sub(time.Now())
		if rem > tokenLeaseTime {
			return errLeaseTimeExceeded
		}
	}
	ch := make(chan struct{})
	var err error
	go func() {
		err = rl.acquireLock(ctx)
		close(ch)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return err
	}
}
func (rl *redisLocker) Release() error {
	pipeliner := rl.client.Pipeline()
	v, err := pipeliner.SRem(tokenSet, rl.lockID).Result()
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}
	if v == 0 {
		return fmt.Errorf("lock has been released")
	}
	err = pipeliner.Publish(rl.lockID, "RELEASED").Err()
	if err != nil {
		return fmt.Errorf("failed to publish lock release: %v", err)
	}
	_, err = pipeliner.Exec()
	if err != nil {
		return fmt.Errorf("failed to release lock: %v", err)
	}
	return nil
}
func (rl *redisLocker) Released() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	<-rl.channel
	return ch
}

func releaseLock(locker *lockerPayload) {
	pipeliner := locker.client.Pipeline()
	pipeliner.SRem(tokenSet, locker.lockID)
	pipeliner.Publish(locker.lockID, "RELEASED")
	pipeliner.Exec()
}

func main() {
	client := conn.NewRedisClient(&conn.RedisOptions{
		Address: "localhost",
		Port:    "6379",
	})

	lockID := "x"

	locker := &lockerPayload{
		lockID:  lockID,
		client:  client,
		channel: client.Subscribe(lockID).Channel(),
	}

	releaseLock(locker)

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()

			// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if taskID == 3 {
				cancel()
			}

			err := acquire(ctx, locker)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer releaseLock(locker)

			// do sth
			dur := time.Duration(rand.Intn(5)+1) * time.Second
			fmt.Printf("task %d: work done in %d seconds\n", taskID, dur/1000000000)
			time.Sleep(dur)
			fmt.Printf("task %d finished\n", taskID)
		}(i)
	}

	wg.Wait()
}
