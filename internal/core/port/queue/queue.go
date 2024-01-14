package queue

import (
	"github.com/kaitsubaka/clubhub_franchises/internal/core/common/queue/event"
)

type Publisher interface {
	Publish(e event.Event) error
}

type Getter interface {
	Get() (event.Event, error)
}

type Subscriber interface {
	Subscribe(e event.Event) error
}
