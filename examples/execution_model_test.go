// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package examples

import (
	"gospec"
	"strings"
)


func ExecutionModelSpec(c *gospec.Context) {
	
	// "Before block", for example common variables for use in all specs.
	commonVariable := ""
	
	c.Specify("The following child specs modify the same variable", func() {
	
		// "Before block", for example initialization for this group of specs.
		commonVariable += "x"
		
		// All sibling specs (specs which are declared within a common parent)
		// are fully isolated from each other. The following three siblings are
		// executed concurrently, each in its own goroutine, and each of them
		// has its own copy of the local variables declared in its parent specs.
		c.Specify("I modify it, but none of my siblings will know it", func() {
			commonVariable += "1"
		})
		c.Specify("Also I modify it, but none of my siblings will know it", func() {
			commonVariable += "2"
		})
		c.Specify("Even I modify it, but none of my siblings will know it", func() {
			commonVariable += "3"
		})
		
		// "After block", for example tear down of changes to the file system.
		commonVariable += "y"
		
		// Depending on which of the previous siblings was executed this time,
		// there are three possible values for the variable:
		c.Then(commonVariable).Should.Be(commonVariable == "x1y" ||
		                                 commonVariable == "x2y" ||
		                                 commonVariable == "x3y")
	})
	
	c.Specify("You can nest", func() {
		c.Specify("as many specs", func() {
			c.Specify("as you wish.", func() {
				c.Specify("GoSpec does not impose artificial limits, " +
				          "so you can organize your specs freely.", func() {
				})
			})
		})
	})
	
	c.Specify("The distinction between 'Should' and 'Must'", func() {
		// When we have non-trivial test setup code, then it is often useful to
		// make assertions about the state of the system under test, before the
		// body of the test is executed. Otherwise it could happen that the test
		// passes even though the code is broken, or then we get lots of
		// unhelpful error messages from the body of the test, even though the
		// bug was in the test setup.
		//
		// For this use case, GoSpec provides a 'Must' in addition to 'Should'.
		// When making assertions about the test setup (i.e. behaviour which is
		// not the focus of the current test) it's better to use 'Must':
		//
		// - When a 'Should' fails, then the child specs are executed normally.
		//
		// - When a 'Must' fails, then the child specs are NOT executed. This
		//   helps to prevent lots of false alarms from the child specs, when
		//   the real problem was in the test setup.
		
		// Some complex test setup code
		input := "abc"
		
		// Uncomment this line to add a bug into the test setup:
		//input += " bug"
		
		// Uncomment one of the following asserts to see their difference:
		//c.Then(input).Should.Equal("abc")
		//c.Then(input).Must.Equal("abc")
		
		c.Specify("When a string is made all uppercase", func() {
			result := strings.ToUpper(input)
			
			c.Specify("Then it is all uppercase", func() {
				c.Then(result).Should.Equal("ABC")
			})
			c.Specify("Its length is not changed", func() {
				c.Then(len(result)).Should.Equal(3)
			})
		})
	})
}

