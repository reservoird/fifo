package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/reservoird/icd"
)

type fifo struct {
	data   *list.List
	mutex  sync.Mutex
	closed bool
	Tag    string
}

// NewQueue is what reservoird to create a queue
func NewQueue() (icd.Queue, error) {
	return new(fifo), nil
}

func (o *fifo) Config(cfg string) error {
	o.data = list.New()
	o.closed = false
	o.Tag = "fifo"
	if cfg != "" {
		d, err := ioutil.ReadFile(cfg)
		if err != nil {
			return err
		}
		f := fifo{}
		err = json.Unmarshal(d, &f)
		if err != nil {
			return err
		}
		o.Tag = f.Tag
	}
	return nil
}

func (o *fifo) Name() string {
	return o.Tag
}

// Put sends data to queue
func (o *fifo) Put(item interface{}) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.closed == false {
		o.data.PushBack(item)
	}
	return fmt.Errorf("fifo is closed")
}

// Get receives data from queue
func (o *fifo) Get() (interface{}, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.closed == false {
		item := o.data.Front()
		if item == nil {
			return nil, fmt.Errorf("fifo is empty")
		}
		value := o.data.Remove(item)
		return value, nil
	}
	return nil, fmt.Errorf("fifo is closed")
}

// Peek receives data from queue
func (o *fifo) Peek() (interface{}, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.closed == false {
		item := o.data.Front()
		if item == nil {
			return nil, fmt.Errorf("fifo is empty")
		}
		return item.Value, nil
	}
	return nil, fmt.Errorf("fifo is closed")
}

// Len returns the current length of the Queue
func (o *fifo) Len() int {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return o.data.Len()
}

// Cap returns the current length of the Queue
func (o *fifo) Cap() int {
	return -1
}

// Clear clears the Queue
func (o *fifo) Clear() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.data.Init()
}

// Closed returns where or not the queue is closed
func (o *fifo) Closed() bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return o.closed
}

// Close closes the channel
func (o *fifo) Close() error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.closed = true
	o.data = nil
	return nil
}
