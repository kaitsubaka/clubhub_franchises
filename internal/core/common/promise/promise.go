package promise

import "sync"

var wg sync.WaitGroup

func All(done chan struct{}, promises ...<-chan any) chan any {
	m := make(chan any, len(promises))
	wg.Add(len(promises))
	for _, c := range promises {
		go multiplex(c, m, done)
	}
	go func() {
		wg.Wait()
		close(m)
	}()
	return m
}

func multiplex(stream <-chan any, multiplexed chan any, done <-chan struct{}) {
	defer wg.Done()
	for v := range stream {
		select {
		case <-done:
			return
		case multiplexed <- v:
		}
	}
}
