package main

import (
	"reflect"
)

type Event interface {
	Yo()
}

type EventA struct {
	Id int
}

type EventB struct {
	Name string
}

func (a EventA) Yo() {
}

func (b EventB) Yo() {
}

type Listener interface {
	GetEvent() Event
	Notify(event Event)
}

type ListenerImpl struct {
	event    Event
	callback func(event Event)
}

func (s *ListenerImpl) GetEvent() Event {
	return s.event
}

func (s *ListenerImpl) Notify(event Event) {
	s.callback(event)
}

type Dispatcher struct {
	backgroundChan chan Event
	stopChan       chan bool
	listeners      map[reflect.Type][]Listener
}

func (s *Dispatcher) Subscribe(callback func(event Event), event Event) *Dispatcher {
	t := reflect.TypeOf(event)
	l, ok := s.listeners[t]
	if !ok {
		l = make([]Listener, 0)
	}
	l = append(l, &ListenerImpl{callback: callback, event: event})
	s.listeners[t] = l
	return s
}

func (s *Dispatcher) Listen() *Dispatcher {
	s.backgroundChan = make(chan Event, 1)
	s.stopChan = make(chan bool, 1)
	go func() {
		for {
			select {
			case <-s.stopChan:
				close(s.stopChan)
				close(s.backgroundChan)
				return
			case e := <-s.backgroundChan:
				t := reflect.TypeOf(e)
				listeners, ok := s.listeners[t]
				if ok {
					for _, l := range listeners {
						l.Notify(e)
					}
				}
			}
		}
	}()
	return s
}

func (s *Dispatcher) Stop() *Dispatcher {
	s.stopChan <- true
	return s
}

func (s *Dispatcher) Fire(event Event) {
	s.backgroundChan <- event
}
