// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"nanospec"
)


func ContextSpec(c nanospec.Context) {

	c.Specify("When specs are executed", func() {
		result := runSpecWithContext(DummySpecWithTwoChildren, newInitialContext())
		executed := result.executedSpecs
		postponed := result.postponedSpecs

		c.Specify("executed specs are reported", func() {
			c.Expect(len(executed)).Equals(2)
			c.Expect(executed[0].name).Equals("RootSpec")
			c.Expect(executed[1].name).Equals("Child A")
		})
		c.Specify("postponed specs are reported", func() {
			c.Expect(len(postponed)).Equals(1)
			c.Expect(postponed[0].name).Equals("Child B")
		})
	})

	c.Specify("When some of the specs have previously been executed", func() {
		result := runSpecWithContext(DummySpecWithTwoChildren, newExplicitContext([]int{1}))
		executed := result.executedSpecs
		postponed := result.postponedSpecs

		c.Specify("previously executed specs are NOT reported", func() {
			c.Expect(len(executed)).Equals(2)
			c.Expect(executed[0].name).Equals("RootSpec")
			c.Expect(executed[1].name).Equals("Child B")

			c.Expect(len(postponed)).Equals(0)
		})
	})

	c.Specify("Postponed specs are scheduled for execution, until they all have been executed", func() {
		r := NewRunner()
		r.AddSpec(DummySpecWithTwoChildren)
		r.Run()

		runCounts := countSpecNames(r.executed)
		c.Expect(len(runCounts)).Equals(3)
		c.Expect(runCounts["gospec.DummySpecWithTwoChildren"]).Equals(2)
		c.Expect(runCounts["Child A"]).Equals(1)
		c.Expect(runCounts["Child B"]).Equals(1)
	})

	c.Specify("Multiple specs can be executed in one batch", func() {
		r := NewRunner()
		r.AddSpec(DummySpecWithOneChild)
		r.AddSpec(DummySpecWithTwoChildren)
		r.Run()

		runCounts := countSpecNames(r.executed)
		c.Expect(runCounts["gospec.DummySpecWithOneChild"]).Equals(1)
		c.Expect(runCounts["gospec.DummySpecWithTwoChildren"]).Equals(2)
	})
}
