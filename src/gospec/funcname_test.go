// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"nanospec"
)


func FuncNameSpec(c nanospec.Context) {

	c.Specify("The name of a function can be retrieved from a function reference", func() {
		name := functionName(dummyFunction)
		c.Expect(name).Equals("gospec.dummyFunction")
	})
	c.Specify("Getting the name of an anonymous functions will fail gracefully", func() {
		name := functionName(func() {})
		c.Expect(name).Equals("<unknown function>")
	})
}

func dummyFunction() {
}
