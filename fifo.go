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

// New is what reservoird to create a queue
func New(cfg string) (icd.Queue, error) {
	o := &fifo{
		data:   list.New(),
		closed: false,
		Tag:    "fifo",
	}
	if cfg != "" {
		d, err := ioutil.ReadFile(cfg)
		if err != nil {
			return nil, err
		}
		f := fifo{}
		err = json.Unmarshal(d, &f)
		if err != nil {
			return nil, err
		}
		o.Tag = f.Tag
	}
	return o, nil
}

func (o *fifo) Name() string {
	return o.Tag
}

// Put sends data to queue
func (o *fifo) Put(item interface{}) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.closed == true {
		return fmt.Errorf("fifo is closed")
	}
	o.data.PushBack(item)
	return nil
}

// Get receives data from queue
func (o *fifo) Get() (interface{}, error) {
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
	return value, nil
}

// Peek receives data from queue
func (o *fifo) Peek() (interface{}, error) {
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
