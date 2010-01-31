// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package examples

import (
	"gospec"
)

func FibSpec(c *gospec.Context) {
	fib := NewFib().Sequence(10)
	
	c.Specify("The first two Fibonacci numbers are 0 and 1", func() {
		c.Then(fib[0]).Should.Equal(0)
		c.Then(fib[1]).Should.Equal(1)
	})
	c.Specify("Each remaining number is the sum of the previous two", func() {
		for i := 2; i < len(fib); i++ {
			c.Then(fib[i]).Should.Equal(fib[i-1] + fib[i-2])
		}
	})
}

