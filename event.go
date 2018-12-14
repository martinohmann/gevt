package gevt

// EventData type definition
type EventData map[string]interface{}

// Event type definition
type Event struct {
	tag  string
	data map[string]interface{}
}

// NewEvent creates a new value of type Event
func NewEvent(tag string, data EventData) Event {
	return Event{
		tag:  tag,
		data: data,
	}
}

// Tag returns the tag of the event
func (e *Event) Tag() string {
	return e.tag
}

// Data returns the event data
func (e *Event) Data() EventData {
	return e.data
}

// Get returns the data for given key
func (e *Event) Get(key string) interface{} {
	return e.data[key]
}

// Has returns true if the key exists in the event data
func (e *Event) Has(key string) bool {
	_, ok := e.data[key]

	return ok
}

// Set sets the value for a key in the event data
func (e *Event) Set(key string, val interface{}) {
	if e.data == nil {
		e.data = EventData{}
	}

	e.data[key] = val
}
