package main

import (
	"fmt"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(4*time.Second),
		sig(2*time.Second),
		sig(3*time.Second),
		sig(3*time.Second),
		sig(1*time.Hour),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	default:
		orDone := make(chan interface{})
		go func() {
			defer close(orDone)
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-or(append(channels[2:], orDone)...):
			}
		}()
		return orDone
	}
}
