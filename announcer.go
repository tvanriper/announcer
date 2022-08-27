package announcer

import (
	"fmt"
	"sync"
)

// Announcer provides the structure with which one may listen to announcements.
type Announcer struct {
	m         sync.Mutex
	listeners map[int]chan<- interface{}
	nextId    int
	capacity  int
	closed    bool
}

// Listener provides the structure with which one may hear announcements from
// the Announcer.
type Listener struct {
	ch        <-chan interface{}
	announcer *Announcer
	id        int
}

// New creates a new Announcer.  Use capacity to establish the amount of
// buffering each Listener should have.
func New(capacity int) *Announcer {
	return &Announcer{capacity: capacity}
}

// Listen creates a new Listener for this Announcer.  Use this to listen for
// announcements.
func (a *Announcer) Listen() *Listener {
	a.m.Lock()
	defer a.m.Unlock()
	if a.listeners == nil {
		a.listeners = make(map[int]chan<- interface{})
	}
	for a.listeners[a.nextId] != nil {
		a.nextId++
	}
	ch := make(chan interface{}, a.capacity)
	if a.closed {
		close(ch)
	}
	a.listeners[a.nextId] = ch
	return &Listener{ch, a, a.nextId}
}

// Send allows an Announcer to send an announcement to listeners.
func (a *Announcer) Send(announcement interface{}) (err error) {
	a.m.Lock()
	defer a.m.Unlock()
	if a.closed {
		return fmt.Errorf("called send after close")
	}
	for _, l := range a.listeners {
		l <- announcement
	}
	return nil
}

func (l *Listener) Listen() <-chan interface{} {
	return l.ch
}

// Close indicates that this listener should receive no more announcements.
func (l *Listener) Close() {
	l.announcer.m.Lock()
	defer l.announcer.m.Unlock()
	delete(l.announcer.listeners, l.id)
}

// Close ends any further announcements from this Announcer.
func (a *Announcer) Close() {
	a.m.Lock()
	defer a.m.Unlock()
	a.closed = true
	for _, l := range a.listeners {
		close(l)
	}
}
