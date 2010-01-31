// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package examples

import (
	"gospec"
)


// Each well-written test has three parts: Arrange, Act, Assert.
// In BDD vocabulary they are often identified by the words Given, When, Then.
//
// A useful style of organizing tests is to have the Arrange and Act in
// the parent spec(s), and then multiple Asserts in the child specs. Because
// GoSpec provides full isolation of common variables in the specs, you may
// organize the specs freely without being worried about side-effects (in the
// following example, see how stack.Pop() is called many times, without any
// order-of-run dependency).
//
// Organizing the tests like this follows the "one assert per test" principle.
// When each test tests only one behaviour, it makes the reason for a test
// failure obvious. When a test fails, you will know exactly what is wrong
// (it isolates the reason for failure, see http://agileinaflash.blogspot.com/2009/02/first.html)
// and you will know whether the behaviour specified by the test is still
// needed, or whether it is obsolete and the test should be removed.
//
// Quite often I use the words Given, When and Then in the test names, because
// they are part of BDD's ubiquitous language. But I always put more emphasis
// on making the tests readable and choosing the best possible words. So when
// it is obvious from the sentence, I may choose to
//   - omit the Given/When/Then keywords,
//   - group the Given and When parts together,
//   - group the When and Then parts together, or even
//   - group all three parts together.
//
// (In the following example, I have marked with comments that which of the
// specs is technically a Give, When or Then. As you can see, there is a
// distinct structure, but also much flexibility.)
//
// For a more thorough example of my style, see
// http://github.com/orfjackal/tdd-tetris-tutorial

func StackSpec(c *gospec.Context) {
	stack := NewStack()
	
	c.Specify("An empty stack", func() { // Given
		
		c.Specify("is empty", func() { // Then
			c.Then(stack).Should.Be(stack.Empty())
		})
		c.Specify("After a push, the stack is no longer empty", func() { // When, Then
			stack.Push("foo")
			c.Then(stack).ShouldNot.Be(stack.Empty())
		})
	})
	
	c.Specify("When objects have been pushed onto a stack", func() { // Given, (When)
		stack.Push("one")
		stack.Push("two")
		
		c.Specify("the object pushed last is popped first", func() { // (When), Then
			x := stack.Pop()
			c.Then(x).Should.Equal("two")
		})
		c.Specify("the object pushed first is popped last", func() { // (When), Then
			stack.Pop()
			x := stack.Pop()
			c.Then(x).Should.Equal("one")
		})
		c.Specify("After popping all objects, the stack is empty", func() { // When, Then
			stack.Pop()
			stack.Pop()
			c.Then(stack).Should.Be(stack.Empty())
		})
	})
}

