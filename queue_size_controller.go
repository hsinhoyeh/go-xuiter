package goxuiter

import (
	"sync"

	log "github.com/golang/glog"
)

type SizeController interface {
	Pause()
	Resume()
	CheckAndWait()
}

type NoOpController struct{}

func (n *NoOpController) Pause()        {}
func (n *NoOpController) Resume()       {}
func (n *NoOpController) CheckAndWait() {}

type QueueSizeController struct {
	sync.Mutex
	wait       chan struct{}
	canEnqueue bool
}

func NewQueueSizeController() *QueueSizeController {
	return &QueueSizeController{
		wait:       make(chan struct{}),
		canEnqueue: true,
	}
}

func (q *QueueSizeController) Pause() {
	if !q.canEnqueue {
		return
	}

	q.Lock()
	q.canEnqueue = false
	q.Unlock()
	log.Infof("sizecontroller: paused")
}

func (q *QueueSizeController) Resume() {
	if q.canEnqueue {
		return
	}

	q.Lock()
	q.canEnqueue = true
	q.Unlock()

	//signal wait chan
	select {
	case q.wait <- struct{}{}:
	default:
	}
	log.Infof("sizecontroller: resumed")
}

func (q *QueueSizeController) CheckAndWait() {
	if !q.canEnqueue {
		<-q.wait
	}
}
