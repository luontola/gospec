// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package examples

import (
	"container/list"
	"gospec"
)


func AssertExamplesSpec(c *gospec.Context) {
	
	c.Specify("Primitives can be compared for equality", func() {
		c.Then(1).Should.Equal(1)
		c.Then(true).Should.Equal(true)
		c.Then("string").Should.Equal("string")
		
		// Comparing floats for equality is not recommended, because
		// the results are not what you might expect. So don't write
		// like this:
		c.Then(3.141).Should.Equal(3.141)
		// But instead write:
		// TODO: c.Then(3.141).Should.BeNear(3.141, 0.001)
	})
	
	c.Specify("Objects with an \"Equals(interface{}) bool\" method can be compared for equality", func() {
		a1 := Point2{1, 2}
		a2 := Point2{1, 2}
		c.Then(a1).Should.Equal(a2)
		
		b1 := &Point3{1, 2, 3}
		b2 := &Point3{1, 2, 3}
		c.Then(b1).Should.Equal(b2)
	})
	
	c.Specify("All assertions can be negated", func() {
		c.Then(1).ShouldNot.Equal(2)
		c.Then("apples").ShouldNot.Equal("oranges")
	})
	
	c.Specify("Any boolean expression can be stated about an object", func() {
		s := "a medium sized string"
		mediumSized := len(s) > 15 && len(s) < 30
		empty := len(s) == 0
		c.Then(s).Should.Be(mediumSized)
		c.Then(s).ShouldNot.Be(empty)
	})
	
	c.Specify("Arrays, slices and containers can be tested for containment", func() {
		array := []string{"one", "two", "three"}
		c.Then(array).Should.Contain("two")
		c.Then(array).ShouldNot.Contain("four")
		
		list := list.New()
		list.PushBack("one")
		list.PushBack("two")
		list.PushBack("three")
		c.Then(list.Iter()).Should.Contain("two")
		c.Then(list.Iter()).ShouldNot.Contain("four")
	})
}

