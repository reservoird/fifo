package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/reservoird/icd"
)

type fifo struct {
	Tag string
}

// NewQueue is what reservoird to create a queue
func NewQueue() (icd.Queue, error) {
	return new(fifo), nil
}

func (o *fifo) Config(cfg string) error {
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
func (o *fifo) Put(data interface{}) error {
	return nil
}

// Get receives data from queue
func (o *fifo) Get() (interface{}, error) {
	return nil, nil
}

// Peek receives data from queue
func (o *fifo) Peek() (interface{}, error) {
	return nil, nil
}

// Len returns the current length of the Queue
func (o *fifo) Len() int {
	return 0
}

// Cap returns the current length of the Queue
func (o *fifo) Cap() int {
	return 0
}

// Clear clears the Queue
func (o *fifo) Clear() {
}

// Closed returns where or not the queue is closed
func (o *fifo) Closed() bool {
	return false
}

// Close closes the channel
func (o *fifo) Close() error {
	return nil
}
