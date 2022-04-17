package priqueue

import (
	"context"
	"sync"

	"github.com/bradenaw/juniper/container/deque"
	"github.com/bradenaw/juniper/stream"
)

type PriorityQueue[P any] struct {
	mw        sync.RWMutex
	eventChan chan interface{}

	q *deque.Deque[P]

	outstream *stream.Stream[P]
}

func NewPriorityQueue[P any](bufferSize int) *PriorityQueue[P] {
	q := &PriorityQueue[P]{
		eventChan: make(chan interface{}, 0),
	}
	return q
}

func (pq *PriorityQueue[P]) PushBack(e P) {

	pq.mw.Lock()
	pq.q.PushBack(e)
	pq.mw.Unlock()
	pq.eventChan <- nil
}

func (pq *PriorityQueue[P]) PushFront(e P) {
	pq.mw.Lock()
	pq.q.PushFront(e)
	pq.mw.Unlock()
	pq.eventChan <- nil
}

func (pq *PriorityQueue[P]) Listen() *stream.Stream[P] {
	if pq.outstream != nil {
		return pq.outstream
	}

	sender, reciver := stream.Pipe[P](0)
	pq.outstream = &reciver

	go func() {
		ctx := context.Background()

		for {
			pq.mw.RLock()
			if pq.q.Len() != 0 {
				err := sender.Send(ctx, pq.q.PopFront())
				if err != nil {
					sender.Close(err)
					reciver = nil
					break
				}
			}
			pq.mw.RUnlock()
		}
	}()

	return &reciver
}
