gevt
====

[![Build Status](https://travis-ci.org/martinohmann/gevt.svg)](https://travis-ci.org/martinohmann/gevt)
[![codecov](https://codecov.io/gh/martinohmann/gevt/branch/master/graph/badge.svg)](https://codecov.io/gh/martinohmann/gevt)
[![Go Report Card](https://goreportcard.com/badge/github.com/martinohmann/gevt)](https://goreportcard.com/report/github.com/martinohmann/gevt)
[![GoDoc](https://godoc.org/github.com/martinohmann/gevt?status.svg)](https://godoc.org/github.com/martinohmann/gevt)

Dead simple event dispatcher. Supports both sync and async dispatching of events.

Installation
------------

```sh
go get -u github.com/martinohmann/gevt
```

Usage
-----

```go
handler := gevt.HandlerFunc(func(event gevt.Event) {
    fmt.Printf("received event %s with data %+v", event.Tag(), event.Data())
})

gevt.Subscribe("foo", handler)

evt := gevt.NewEvent("foo", gevt.EventData{"bar": 1})

<-gevt.Publish(evt)

gevt.Unsubscribe("foo", handler)
```

Check [example_test.go](example_test.go) or
[godoc](https://godoc.org/github.com/martinohmann/gevt) for more usage
examples.

License
-------

The source code of this is released under the MIT License. See the bundled LICENSE
file for details.
