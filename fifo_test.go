package main

import (
	"testing"

	"github.com/reservoird/icd"
)

func TestChanQueueImplements(t *testing.T) {
	q, err := NewQueue()
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	err = q.Config("")
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	i := interface{}(q)
	_, ok := i.(icd.Queue)
	if ok == false {
		t.Errorf("error queue does not implement icd.Queue")
	}
}

func TestChanQueueConfig(t *testing.T) {
	q, err := NewQueue()
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	err = q.Config("")
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
}

func TestChanQueuePut(t *testing.T) {
	q, err := NewQueue()
	if err != nil {
		t.Errorf("expecting nil got error %v", err)
	}
	err = q.Config("")
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	s := "hello"
	d := []byte(s)
	err = q.Put(d)
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
}

func TestChanQueuePutErrorNil(t *testing.T) {
	q, err := NewQueue()
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	err = q.Config("")
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	d := []byte(nil)
	err = q.Put(d)
	if err == nil {
		t.Errorf("error expecting error got nil")
	}
}

func TestChanQueuePutErrorClosed(t *testing.T) {
	q, err := NewQueue()
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	err = q.Config("")
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	d := []byte(nil)
	err = q.Close()
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	err = q.Put(d)
	if err == nil {
		t.Errorf("error expecting error got nil")
	}
}

func TestChanQueueGet(t *testing.T) {
	q, err := NewQueue()
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	err = q.Config("")
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	s := "hello"
	d := []byte(s)
	err = q.Put(d)
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	e, err := q.Get()
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	actual, ok := e.([]byte)
	if ok == false {
		t.Errorf("error invalid type")
	}
	sa := string(actual)
	if s != sa {
		t.Errorf("error expecting %s got %s", s, sa)
	}
}

func TestChanQueuePopErrorClosed(t *testing.T) {
	q, err := NewQueue()
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	err = q.Config("")
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	s := "hello"
	d := []byte(s)
	err = q.Put(d)
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	err = q.Close()
	if err != nil {
		t.Errorf("error expecting nil go error %v", err)
	}
	_, err = q.Get()
	if err == nil {
		t.Errorf("error expecting error got nil")
	}
}

func TestChanQueueClose(t *testing.T) {
	q, err := NewQueue()
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	err = q.Config("")
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	err = q.Close()
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
}
