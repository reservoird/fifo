package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

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
		cfg:    c,
		data:   list.New(),
		mutex:  sync.Mutex{},
		stats:  FifoStats{},
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

// Peek receives data from queue
func (o *Fifo) Peek() (interface{}, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.closed == true {
		return nil, fmt.Errorf("fifo is closed")
	}
	item := o.data.Front()
	if item == nil {
		return nil, nil
	}
	return item.Value, nil
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
	return nil
}

func (o *Fifo) getStats(monitoring bool) (string, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.stats.Name = o.cfg.Name
	o.stats.Monitoring = monitoring

	data, err := json.Marshal(o.stats)
	if err != nil {
		return "", err
	}
	return string(data), nil
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
func (o *Fifo) Monitor(statsChan chan<- string, clearChan <-chan struct{}, doneChan <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done() // required

	run := true
	for run == true {
		// clear
		select {
		case <-clearChan:
			o.clearStats()
		default:
		}

		// done
		select {
		case <-doneChan:
			run = false
		default:
		}

		// get stats
		stats, err := o.getStats(run)
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			// stats
			select {
			case statsChan <- stats:
			default:
			}
		}

		if run == true {
			time.Sleep(time.Millisecond)
		}
	}
}
