// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package examples

import (
	"github.com/orfjackal/gospec/src/gospec"
	. "github.com/orfjackal/gospec/src/gospec"
)

// This is the style that I've found the most useful in organizing tests.
// In the parent spec(s) is done some action and then it has multiple child
// specs which each verify one isolated piece of behaviour. Each spec has a
// well though out name which explains the motivation behind the code.
//
// To learn more, see this article and tutorial:
// http://blog.orfjackal.net/2010/02/three-styles-of-naming-tests.html
// http://github.com/orfjackal/tdd-tetris-tutorial

func StackSpec(c gospec.Context) {
	stack := NewStack()

	c.Specify("An empty stack", func() {

		c.Specify("is empty", func() {
			c.Expect(stack.Empty(), IsTrue)
		})
		c.Specify("After a push, the stack is no longer empty", func() {
			stack.Push("a push")
			c.Expect(stack.Empty(), IsFalse)
		})
	})

	c.Specify("When objects have been pushed onto a stack", func() {
		stack.Push("pushed first")
		stack.Push("pushed last")

		c.Specify("the object pushed last is popped first", func() {
			poppedFirst := stack.Pop()
			c.Expect(poppedFirst, Equals, "pushed last")
		})
		c.Specify("the object pushed first is popped last", func() {
			stack.Pop()
			poppedLast := stack.Pop()
			c.Expect(poppedLast, Equals, "pushed first")
		})
		c.Specify("After popping all objects, the stack is empty", func() {
			stack.Pop()
			stack.Pop()
			c.Expect(stack.Empty(), IsTrue)
		})
	})
}
