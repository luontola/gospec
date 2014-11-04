// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"github.com/orfjackal/nanospec.go/src/nanospec"
	"sort"
)

func ExecutionModelSpec(c nanospec.Context) {

	c.Specify("Specs with children, but without siblings, are executed fully on one run", func() {

		c.Specify("Case: no children", func() {
			runSpecWithContext(DummySpecWithNoChildren, newInitialContext())
			c.Expect(testSpy).Equals("root")
		})
		c.Specify("Case: one child", func() {
			runSpecWithContext(DummySpecWithOneChild, newInitialContext())
			c.Expect(testSpy).Equals("root,a")
		})
		c.Specify("Case: nested children", func() {
			runSpecWithContext(DummySpecWithNestedChildren, newInitialContext())
			c.Expect(testSpy).Equals("root,a,aa")
		})
	})

	c.Specify("Specs with siblings are executed only one sibling at a time", func() {

		c.Specify("Case: on initial run, the 1st child is executed", func() {
			runSpecWithContext(DummySpecWithTwoChildren, newInitialContext())
			c.Expect(testSpy).Equals("root,a")
		})
		c.Specify("Case: explicitly execute the 1st child", func() {
			runSpecWithContext(DummySpecWithTwoChildren, newExplicitContext([]int{0}))
			c.Expect(testSpy).Equals("root,a")
		})
		c.Specify("Case: explicitly execute the 2nd child", func() {
			runSpecWithContext(DummySpecWithTwoChildren, newExplicitContext([]int{1}))
			c.Expect(testSpy).Equals("root,b")
		})
	})

	c.Specify("Specs with nested siblings: eventually all siblings are executed, one at a time, in isolation", func() {
		r := NewParallelRunner()
		r.AddSpec(DummySpecWithMultipleNestedChildren)

		// Execute manually instead of calling Run(), in order to avoid running
		// the specs multi-threadedly, which would mess up the test spy.
		runs := make([]string, 0)
		for r.hasScheduledTasks() {
			resetTestSpy()
			r.executeNextScheduledTask()
			runs = append(runs, testSpy)
		}
		sort.Strings(runs)

		c.Expect(len(runs)).Equals(5)
		c.Expect(runs[0]).Equals("root,a,aa")
		c.Expect(runs[1]).Equals("root,a,ab")
		c.Expect(runs[2]).Equals("root,b,ba")
		c.Expect(runs[3]).Equals("root,b,bb")
		c.Expect(runs[4]).Equals("root,b,bc")
	})
}
