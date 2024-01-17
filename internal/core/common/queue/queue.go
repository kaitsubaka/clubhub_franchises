package queue

import (
	"sync"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/common/queue/event"
	qport "github.com/kaitsubaka/clubhub_franchises/internal/core/port/queue"
)

type Queue struct {
	q          chan event.Event
	wg         sync.WaitGroup
	done       chan struct{}
	numWorkers uint
	subscriber qport.Subscriber
}

func New(capacity uint) *Queue {
	return &Queue{
		q:          make(chan event.Event, capacity),
		done:       make(chan struct{}),
		numWorkers: 1,
	}
}

func (q *Queue) WithSubscriber(s qport.Subscriber) *Queue {
	q.subscriber = s
	return q
}

func (q *Queue) WithNumOfWorker(n uint) *Queue {
	q.numWorkers = n
	return q
}

func (q *Queue) Build() (*Queue, error) {
	for i := 0; i < int(q.numWorkers); i++ {
		go func(q *Queue) {
			defer q.wg.Done()
		Outer:
			for {
				select {
				case e, ok := <-q.q:
					if !ok && len(q.q) == 0 {
						break Outer
					}
					q.subscriber.Subscribe(e)
				}
			}
		}(q)
	}
	q.wg.Add(int(q.numWorkers))
	go func() {
		q.wg.Wait()
		q.done <- struct{}{}
	}()
	return q, nil
}

func (q *Queue) Close() error {
	close(q.q)
	<-q.done
	return nil
}

func (q *Queue) Publish(e event.Event) error {
	q.q <- e
	return nil
}

func (q *Queue) Get() (event.Event, error) {
	return <-q.q, nil
}
