package gevt

import (
	"sync/atomic"
	"testing"
)

type testHandler struct {
	called int64
}

func (h *testHandler) HandleEvent(event Event) {
	atomic.AddInt64(&h.called, 1)
}

func TestPublishSync(t *testing.T) {
	h := &testHandler{}

	evt := NewEvent("foo", nil)

	Subscribe("foo", h)
	PublishSync(evt)

	Unsubscribe("foo", h)
	PublishSync(evt)

	if c := atomic.LoadInt64(&h.called); c != 1 {
		t.Fatalf("expected handler to be called once, was called %d times", c)
	}
}

func TestPublish(t *testing.T) {
	h := &testHandler{}

	evt := NewEvent("foo", nil)

	Subscribe("foo", h)
	one := Publish(evt)
	two := Publish(evt)

	<-two
	<-one

	Unsubscribe("foo", h)

	if c := atomic.LoadInt64(&h.called); c != 2 {
		t.Fatalf("expected handler to be called twice, was called %d times", c)
	}
}

func TestPanickingHandler(t *testing.T) {
	h := HandlerFunc(func(event Event) {
		panic("oops")
	})

	evt := NewEvent("foo", nil)

	Subscribe("foo", h)
	<-Publish(evt)
}

func TestMultiSubscribe(t *testing.T) {
	h := &testHandler{}

	Subscribe("foo", h, h)

	evt := NewEvent("foo", nil)

	<-Publish(evt)

	UnsubscribeAll("foo")

	if c := atomic.LoadInt64(&h.called); c != 1 {
		t.Fatalf("expected handler to be called once, was called %d times", c)
	}
}
