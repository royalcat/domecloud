package entrycache

import (
	"context"
	"sync"
	"time"

	"github.com/bradenaw/juniper/xsync/xatomic"
)

type MediaCacheWorker struct {
	mc *EntryIndex
}

type EventFirer struct {
	minDuration time.Duration

	bumped xatomic.Value[bool]

	queueLock sync.Mutex
	fireQueue []func()
}

func (ef *EventFirer) Bump() {
	ef.bumped.Store(true)
}

func (ef *EventFirer) Add() {
	ef.queueLock.Lock()

	ef.queueLock.Unlock()
}

func (ef *EventFirer) Run(ctx context.Context) {
	ticker := time.NewTicker(ef.minDuration)
	for {
		select {
		case <-ticker.C:
			if ef.bumped.Load() == false {
				for ef.bumped.Load() == false {
					ef.queueLock.Lock()
					if len(ef.fireQueue) != 0 {
						ef.queueLock.Unlock()
						break
					}
					work := ef.fireQueue[0]         // get first event
					ef.fireQueue = ef.fireQueue[1:] // pop runned event
					ef.queueLock.Unlock()
					work() // run first event from queue
				}
			} else {
				ef.bumped.Store(false)
			}
		case <-ctx.Done():
			break
		}
	}
}
