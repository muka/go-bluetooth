package emitter

import (
	"testing"
)

func TestEmitterSimple(t *testing.T) {

	t.Log("On")

	fn := NewCallback(func(ev Event) {
		if ev.GetData() == "Hello World" {
			t.Log("Event received")
		}
	})

	On("test", fn)

	t.Log("Emit")
	Emit("test", "Hello World")

	t.Log("Off")
	Off("test", fn)

}
