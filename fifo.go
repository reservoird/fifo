package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"sync"

	"github.com/reservoird/icd"
)

// FifoCfg contains config
type FifoCfg struct {
	Name string
}

// FifoStats contains stats
type FifoStats struct {
	Name             string
	MessagesReceived uint64
	MessagesSent     uint64
	Len              uint64
	Closed           bool
	Monitoring       bool
}

// Fifo contains what is needed for queue
type Fifo struct {
	cfg    FifoCfg
	data   *list.List
	mutex  sync.Mutex
	stats  FifoStats
	closed bool
}

// New is what reservoird to create a queue
func New(cfg string) (icd.Queue, error) {
	c := FifoCfg{
		Name: "com.reservoird.queue.fifo",
	}
	if cfg != "" {
		d, err := ioutil.ReadFile(cfg)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(d, &c)
		if err != nil {
			return nil, err
		}
	}
	o := &Fifo{
		cfg:   c,
		data:  list.New(),
		mutex: sync.Mutex{},
		stats: FifoStats{
			Name: c.Name,
		},
		closed: false,
	}
	return o, nil
}

// Name returns the name
func (o *Fifo) Name() string {
	return o.cfg.Name
}

// Put sends data to queue
func (o *Fifo) Put(item interface{}) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.stats.MessagesReceived = o.stats.MessagesReceived + 1
	if o.closed == true {
		return fmt.Errorf("fifo is closed")
	}
	o.data.PushBack(item)
	return nil
}

// Get receives data from queue
func (o *Fifo) Get() (interface{}, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.closed == true {
		return nil, fmt.Errorf("fifo is closed")
	}
	item := o.data.Front()
	if item == nil {
		return nil, nil
	}
	value := o.data.Remove(item)
	o.stats.MessagesSent = o.stats.MessagesSent + 1
	return value, nil
}

// Len returns the current length of the Queue
func (o *Fifo) Len() int {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.stats.Len = uint64(o.data.Len())
	return o.data.Len()
}

// Cap returns the current length of the Queue
func (o *Fifo) Cap() int {
	return -1
}

// Clear clears the Queue
func (o *Fifo) Clear() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.data.Init()
}

// Closed returns where or not the queue is closed
func (o *Fifo) Closed() bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return o.closed
}

// Close closes the channel
func (o *Fifo) Close() error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.closed = true
	o.data = nil
	o.stats.Closed = o.closed
	return nil
}

func (o *Fifo) getStats(monitoring bool) *FifoStats {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	fifoStats := &FifoStats{
		Name:             o.cfg.Name,
		MessagesReceived: o.stats.MessagesReceived,
		MessagesSent:     o.stats.MessagesSent,
		Len:              o.stats.Len,
		Closed:           o.stats.Closed,
		Monitoring:       monitoring,
	}

	return fifoStats
}

func (o *Fifo) clearStats() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.stats = FifoStats{
		Name:   o.cfg.Name,
		Closed: o.closed,
	}
}

// Monitor provides statistics and clear
func (o *Fifo) Monitor(mc *icd.MonitorControl) {
	defer mc.WaitGroup.Done() // required

	run := true
	for run == true {
		// clear
		select {
		case <-mc.ClearChan:
			o.clearStats()
		default:
		}

		// done
		select {
		case <-mc.DoneChan:
			run = false
		default:
		}

		// send stats
		select {
		case mc.StatsChan <- o.getStats(run):
		default:
		}

		runtime.Gosched()
	}

	// send final stats blocking
	mc.FinalStatsChan <- o.getStats(run)
}
