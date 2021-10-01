package goxuiter

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
	//"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/queue"
)

type CollyController struct {
	*colly.Collector

	qc                SizeController
	q                 *queue.Queue
	queueSizeCallback []QueueSizeCallback
}

type QueueSizeCallback func(size int)

func NewCollyController(qc SizeController, concurrency int) *CollyController {
	c := colly.NewCollector(
		//colly.Async(),
		colly.AllowURLRevisit(),
		// Attach a debugger to the collector
		//colly.Debugger(&debug.LogDebugger{}),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36"),
	)
	c.WithTransport(&http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost: 6000,
	})

	q, _ := queue.New(
		int(concurrency), // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 1000 * 1000 * 1000 * 1000 /*some value we never cross*/},
	)

	cc := &CollyController{
		Collector: c,
		qc:        qc,
		q:         q,
	}
	go cc.monitorLoop()
	return cc
}

func (c *CollyController) monitorLoop() {
	tick := time.Tick(1 * time.Second)
	for range tick {
		if len(c.queueSizeCallback) == 0 {
			continue
		}
		s, _ := c.q.Size()
		for _, callback := range c.queueSizeCallback {
			callback(s)
		}
	}
}

func (c *CollyController) OnQueueSizeMonitor(f QueueSizeCallback) {
	c.queueSizeCallback = append(c.queueSizeCallback, f)
}

func (c *CollyController) AddURL(url string) error {
	c.qc.CheckAndWait()
	return c.q.AddURL(url)
}

func (c *CollyController) AddRequest(r *colly.Request) error {
	c.qc.CheckAndWait()
	return c.q.AddRequest(r)
}

func (c *CollyController) Run() {
	if err := c.q.Run(c.Collector); err != nil {
		panic(err)
	}
	c.Wait()
}
