// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package examples

import (
	"gospec"
)


// An example of how one might write specs using GoSpec.
func StackSpec(c *gospec.Context) {
	stack := NewStack()
	
	c.Specify("An empty stack", func() {
		
		c.Specify("is empty", func() {
			c.Then(stack).Should.Be(stack.Empty())
		})
		c.Specify("After a push, the stack is no longer empty", func() {
			stack.Push("foo")
			c.Then(stack).ShouldNot.Be(stack.Empty())
		})
	})
	
	c.Specify("When objects have been pushed onto a stack", func() {
		stack.Push("one")
		stack.Push("two")
		
		c.Specify("the object pushed last is popped first", func() {
			x := stack.Pop()
			c.Then(x).Should.Equal("two")
		})
		c.Specify("the object pushed first is popped last", func() {
			stack.Pop()
			x := stack.Pop()
			c.Then(x).Should.Equal("one")
		})
		c.Specify("After popping all objects, the stack is empty", func() {
			stack.Pop()
			stack.Pop()
			c.Then(stack).Should.Be(stack.Empty())
		})
	})
}

