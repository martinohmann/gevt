package gevt_test

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/martinohmann/gevt"
)

func ExampleNewDispatcher() {
	d := gevt.NewDispatcher()

	<-d.Publish(gevt.NewEvent("foo", nil))
}

func ExamplePublish() {
	slowHandler := gevt.HandlerFunc(func(event gevt.Event) {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("slow handler done")
	})

	fastHandler := gevt.HandlerFunc(func(event gevt.Event) {
		fmt.Println("fast handler done")
	})

	gevt.Subscribe("foo", slowHandler, fastHandler)

	evt := gevt.NewEvent("foo", nil)

	<-gevt.Publish(evt)

	gevt.Unsubscribe("foo")

	// Output:
	// fast handler done
	// slow handler done
}

func ExamplePublishSync() {
	handler := gevt.HandlerFunc(func(event gevt.Event) {
		fmt.Println("handler was called synchronously")
	})

	gevt.Subscribe("foo", handler)

	evt := gevt.NewEvent("foo", nil)

	gevt.PublishSync(evt)

	gevt.Unsubscribe("foo", handler)

	// Output:
	// handler was called synchronously
}

type badHandler struct{}

func (*badHandler) HandleEvent(event gevt.Event) {
	panic("oops")
}

func ExampleWithLogger() {
	// type badHandler struct{}

	// func (*badHandler) HandleEvent(event gevt.Event) {
	// 	panic("oops")
	// }

	gevt.WithLogger(log.New(os.Stdout, "gevt: ", 0))

	handler := &badHandler{}

	gevt.Subscribe("foo", handler)

	evt := gevt.NewEvent("foo", nil)

	gevt.PublishSync(evt)

	gevt.Unsubscribe("foo", handler)

	// Output:
	// gevt: handler &gevt_test.badHandler{} panicked on event gevt.Event{tag:"foo", data:map[string]interface {}(nil)} with: oops
}

func ExampleSubscribe() {
	h1 := gevt.HandlerFunc(func(event gevt.Event) {})
	h2 := gevt.HandlerFunc(func(event gevt.Event) {})

	// subscribe h1 to events with tag "foo"
	gevt.Subscribe("foo", h1)

	// subscribe h1 and h2 to events with tag "foo"
	gevt.Subscribe("foo", h1, h2)
}

func ExampleUnsubscribe() {
	h1 := gevt.HandlerFunc(func(event gevt.Event) {})
	h2 := gevt.HandlerFunc(func(event gevt.Event) {})

	// unsubscribe h1 from events with tag "foo"
	gevt.Unsubscribe("foo", h1)

	// unsubscribe h1 and h2 from events with tag "foo"
	gevt.Unsubscribe("foo", h1, h2)

	// unsubscribe all handlers from "foo" event
	gevt.Unsubscribe("foo")
}

func ExampleUnsubscribeAll() {
	// unsubscribe handlers for "foo" tag
	gevt.UnsubscribeAll("foo")

	// unsubscribe handlers for "foo" and "bar" tags
	gevt.UnsubscribeAll("foo", "bar")

	// unsubscribe handlers for all tags
	gevt.UnsubscribeAll()
}

func ExampleEvent() {
	evt := gevt.NewEvent("foo", gevt.EventData{
		"bar": 1,
	})

	fmt.Println(evt.Tag())
	fmt.Println(evt.Has("nope"))
	fmt.Println(evt.Get("bar"))
	fmt.Println(evt.Data())

	// Output:
	// foo
	// false
	// 1
	// map[bar:1]
}
