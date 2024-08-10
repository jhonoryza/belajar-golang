package gocontext

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1

		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
			}
		}
	}()

	return destination
}

func TestCancel(t *testing.T) {
	fmt.Println("total goroutine", runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("counter", n)
		if n == 10 {
			break
		}
	}
	cancel()

	time.Sleep(2 * time.Second)

	fmt.Println("total goroutine", runtime.NumGoroutine())
}

func TestCancelWithTimeout(t *testing.T) {
	fmt.Println("total goroutine", runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 2*time.Second)
	defer cancel()

	destination := CreateSlowCounter(ctx)
	for n := range destination {
		fmt.Println("counter", n)
		if n == 10 {
			break
		}
	}

	fmt.Println("total goroutine", runtime.NumGoroutine())
}

func CreateSlowCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1

		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second) // simulate slow counter
			}
		}
	}()

	return destination
}

func TestCancelWithDeadline(t *testing.T) {
	fmt.Println("total goroutine", runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(2*time.Second))
	defer cancel()

	destination := CreateSlowCounter(ctx)
	for n := range destination {
		fmt.Println("counter", n)
		if n == 10 {
			break
		}
	}

	fmt.Println("total goroutine", runtime.NumGoroutine())
}
