package gevt

import (
	"reflect"
	"testing"
)

func TestEvent(t *testing.T) {
	e := NewEvent("sometag", nil)

	if e.Has("foo") {
		t.Fatalf("expected event not to contain foo")
	}

	res := e.Get("foo")
	if res != nil {
		t.Fatalf("expected res to be <nil>, got %v", res)
	}

	e.Set("foo", 42)

	if !e.Has("foo") {
		t.Fatalf("expected event contain foo")
	}

	res = e.Get("foo")
	if res != 42 {
		t.Fatalf("expected res to be 42, got %v", res)
	}

	data := e.Data()
	if reflect.DeepEqual(data, map[string]interface{}{"foo": 42}) {
		t.Fatalf("expected event data not to be <nil>")
	}
}
