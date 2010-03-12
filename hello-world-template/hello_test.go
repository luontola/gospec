package hello

import (
	"gospec"
	. "gospec"
)


func HelloSpec(c gospec.Context) {

	c.Specify("Says a friendly greeting", func() {
		c.Expect(SayHello("World"), Equals, "Hello, World")
	})
}

