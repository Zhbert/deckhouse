package template

import (
	"sync"
	"time"
)

type changesEmitter struct {
	sync.Mutex
	count int
}

func (e *changesEmitter) emitChanges() {
	e.Lock()
	e.count++
	e.Unlock()
}

func (e *changesEmitter) runBufferedEmitter(channel chan struct{}) {
	for {
		// we need sleep to avoid emitting configuration change on batch updates
		// for example on a start - we add all NodeGroupConfigurations, but need to rerender context and checksums only once
		time.Sleep(500 * time.Millisecond)
		e.Lock()
		if e.count > 0 {
			channel <- struct{}{}
			e.count = 0
		}
		e.Unlock()
	}
}
