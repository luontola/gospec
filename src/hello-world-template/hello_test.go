package hello

import (
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)

func HelloSpec(c gospec.Context) {

	c.Specify("Says a friendly greeting", func() {
		c.Expect(SayHello("World"), Equals, "Hello, World")
	})
}
