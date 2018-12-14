// Package gevt provides a dead simple event dispatcher and interfaces for
// creating event handlers.
package gevt

import (
	"log"
	"sync"
)

// Handler defines the interface for an event handler
type Handler interface {
	// HandleEvent handles an event
	HandleEvent(event Event)
}

// handlerFunc wraps a plain for usage as an event handler
type handlerFunc struct {
	f func(event Event)
}

// HandlerFunc wraps a plain function into an object that safisfies the
// Handler interface
func HandlerFunc(f func(event Event)) *handlerFunc {
	return &handlerFunc{f}
}

// HandleEvent implements the Handler interface
func (h *handlerFunc) HandleEvent(event Event) {
	h.f(event)
}

// Dispatcher defines the interface for an event dispatcher
type Dispatcher interface {
	// WithLogger sets the logger for the dispatcher
	WithLogger(logger *log.Logger)

	// Subscribe subscribes handlers to events with given tag
	Subscribe(tag string, handler ...Handler)

	// Unsubscribe unsubscribes handlers from events with given tag. If no
	// handlers are provided, all handlers for the given tag will be
	// unsubscribed.
	Unsubscribe(tag string, handler ...Handler)

	// UnsubscribeAll unsubscribes all handlers from all events with given
	// tags. If no tags are provided, handlers for all tags will be
	// unsubscribed.
	UnsubscribeAll(tags ...string)

	// PublishSync publishes an event synchronously. This will block until all
	// event handlers processed the event.
	PublishSync(event Event)

	// Publish publishes an event and calls each handler in a separate
	// goroutine. It returns a channel which can be used to get notified when
	// all event handlers have processed the event.
	Publish(event Event) <-chan bool
}

var defaultDispatcher = NewDispatcher()

type dispatcher struct {
	mu          sync.Mutex
	subscribers map[string][]Handler
	logger      *log.Logger
}

// NewDispatcher creates a new event dispatcher
func NewDispatcher() Dispatcher {
	return &dispatcher{
		subscribers: make(map[string][]Handler),
	}
}

// WithLogger sets the logger for the dispatcher
func (d *dispatcher) WithLogger(logger *log.Logger) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.logger = logger
}

// Subscribe subscribes handlers to events with given tag
func (d *dispatcher) Subscribe(tag string, handlers ...Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(handlers) == 0 {
		return
	}

	if _, ok := d.subscribers[tag]; !ok {
		d.subscribers[tag] = []Handler{}
	}

	for _, h := range handlers {
		d.subscribe(tag, h)
	}
}

// subscribe subscribes a handler to events for the given tag
func (d *dispatcher) subscribe(tag string, handler Handler) {
	if handler == nil {
		return
	}

	for _, h := range d.subscribers[tag] {
		if h == handler {
			return
		}
	}

	d.subscribers[tag] = append(d.subscribers[tag], handler)
}

// Unsubscribe unsubscribes handlers from events with given tag. If no
// handlers are provided, all handlers for the given tag will be
// unsubscribed.
func (d *dispatcher) Unsubscribe(tag string, handlers ...Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(handlers) == 0 {
		delete(d.subscribers, tag)
		return
	}

	for _, h := range handlers {
		d.unsubscribe(tag, h)
	}
}

// unsubscribe unsubscribes a handler from events for the given tag
func (d *dispatcher) unsubscribe(tag string, handler Handler) {
	s := d.subscribers[tag]

	for i, h := range s {
		if h == handler {
			d.subscribers[tag] = append(s[:i], s[i+1:]...)
			return
		}
	}
}

// UnsubscribeAll unsubscribes all handlers from all events with given
// tags. If no tags are provided, handlers for all tags will be
// unsubscribed.
func (d *dispatcher) UnsubscribeAll(tags ...string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(tags) == 0 {
		d.subscribers = make(map[string][]Handler)
		return
	}

	for _, tag := range tags {
		delete(d.subscribers, tag)
	}
}

// PublishSync publishes an event synchronously. This will block until all
// event handlers processed the event.
func (d *dispatcher) PublishSync(event Event) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, s := range d.subscribers[event.Tag()] {
		d.publish(s, event)
	}
}

// Publish publishes an event and calls each handler in a separate
// goroutine. It returns a channel which can be used to get notified when
// all event handlers have processed the event.
func (d *dispatcher) Publish(event Event) <-chan bool {
	done := make(chan bool)

	var wg sync.WaitGroup

	d.mu.Lock()
	defer func() {
		d.mu.Unlock()
		wg.Wait()
		close(done)
	}()

	for _, s := range d.subscribers[event.Tag()] {
		wg.Add(1)
		go func(h Handler, event Event) {
			defer wg.Done()
			d.publish(h, event)
		}(s, event)
	}

	return done
}

// publish publishes an event to the handler make sure to catch panics
func (d *dispatcher) publish(handler Handler, event Event) {
	defer func() {
		if err := recover(); err != nil {
			d.log("handler %#v panicked on event %#v with: %s", handler, event, err)
		}
	}()

	handler.HandleEvent(event)
}

// log logs a message if a logger is set
func (d *dispatcher) log(fmt string, v ...interface{}) {
	if d.logger != nil {
		d.logger.Printf(fmt, v...)
	}
}

// WithLogger sets the logger for the default dispatcher
func WithLogger(logger *log.Logger) {
	defaultDispatcher.WithLogger(logger)
}

// Subscribe subscribes handlers to events with given tag on the default event
// dispatcher
func Subscribe(tag string, handlers ...Handler) {
	defaultDispatcher.Subscribe(tag, handlers...)
}

// Unsubscribe unsubscribes handlers from the defualt event dispatcher for
// events with given tag. If no handlers are provided, all handlers for
// the given tag will be unsubscribed.
func Unsubscribe(tag string, handlers ...Handler) {
	defaultDispatcher.Unsubscribe(tag, handlers...)
}

// UnsubscribeAll unsubscribes all handlers from the default event dispatcher
// for all events with given tags. If no tags are provided, handlers for all
// tags will be unsubscribed.
func UnsubscribeAll(tags ...string) {
	defaultDispatcher.UnsubscribeAll(tags...)
}

// PublishSync publishes an event synchronously using the default event
// dispatcher. This will block until all event handlers processed the event.
func PublishSync(event Event) {
	defaultDispatcher.PublishSync(event)
}

// Publish publishes an event using the default event dispatcher and calls each
// handlers in a separete goroutine. It returns a channel which can be used to
// get notified when all event handlers have processed the event.
func Publish(event Event) <-chan bool {
	return defaultDispatcher.Publish(event)
}
