// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"testing"
)


func Test__When_a_spec_has_passing_expectations__Then_the_spec_passes(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Expect(1, Equals, 1)
	})
	assertEquals(1, results.PassCount(), t)
	assertEquals(0, results.FailCount(), t)
}

func Test__When_a_spec_has_failing_expectations__Then_the_spec_fails(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Expect(1, Equals, 2)
	})
	assertEquals(0, results.PassCount(), t)
	assertEquals(1, results.FailCount(), t)
}


func Test__When_a_spec_has_passing_expectations__Then_its_children_are_executed(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Expect(1, Equals, 1)
		c.Specify("Child", func() {
		})
	})
	assertEquals(2, results.TotalCount(), t)
}

func Test__When_a_spec_has_failing_expectations__Then_its_children_are_executed(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Expect(1, Equals, 2)
		c.Specify("Child", func() {
		})
	})
	assertEquals(2, results.TotalCount(), t)
}

func Test__When_a_spec_has_passing_assumptions__Then_its_children_are_executed(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Assume(1, Equals, 1)
		c.Specify("Child", func() {
		})
	})
	assertEquals(2, results.TotalCount(), t)
}

func Test__When_a_spec_has_failing_assumptions__Then_its_children_are_NOT_executed(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Assume(1, Equals, 2)
		c.Specify("Child", func() {
		})
	})
	assertEquals(1, results.TotalCount(), t)
}


// Location

func Test__The_location_of_a_failed_expectation_is_reported(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Expect(1, Equals, 2)
	})
	assertErrorIsInFile("expectations_test.go", results, t)
}

func Test__The_location_of_a_failed_assumption_is_reported(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Assume(1, Equals, 2)
	})
	assertErrorIsInFile("expectations_test.go", results, t)
}

func assertErrorIsInFile(file string, results *ResultCollector, t *testing.T) {
	for spec := range results.sortedRoots() {
		error := spec.errors.Front().Value.(*Error)
		assertEquals(file, error.Location.File, t)
	}
}

