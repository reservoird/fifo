package main

import (
	"testing"

	"github.com/reservoird/icd"
)

func TestFifoImplements(t *testing.T) {
	q, err := New("", nil)
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	i := interface{}(q)
	_, ok := i.(icd.Queue)
	if ok == false {
		t.Errorf("error queue does not implement icd.Queue")
	}
}

func TestFifoConfig(t *testing.T) {
	_, err := New("", nil)
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
}

func TestFifoPut(t *testing.T) {
	q, err := New("", nil)
	if err != nil {
		t.Errorf("expecting nil got error %v", err)
	}
	s := "hello"
	d := []byte(s)
	err = q.Put(d)
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
}

func TestFifoPutErrorClosed(t *testing.T) {
	q, err := New("", nil)
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

func TestFifoGet(t *testing.T) {
	q, err := New("", nil)
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

func TestFifoGetErrorClosed(t *testing.T) {
	q, err := New("", nil)
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

func TestFifoClose(t *testing.T) {
	q, err := New("", nil)
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
	err = q.Close()
	if err != nil {
		t.Errorf("error expecting nil got error %v", err)
	}
}
