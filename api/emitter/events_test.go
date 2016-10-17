package emitter

import (
	"testing"
)

func TestEmitterSimple(t *testing.T) {

	t.Log("On")
	On("test", func(ev Event) {
		if ev.Data == "Hello World" {
			t.Log("Event received")
		}
	})

	t.Log("Emit")
	Emit("test", "Hello World")

	t.Log("Off")
	Off("test")

}
