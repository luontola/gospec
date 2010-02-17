// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"testing"
)


// TODO: remove when the c.Then() syntax is removed

func Test__When_a_spec_contains_passing_asserts__Then_the_spec_passes(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Then(1).Should.Equal(1)
	})
	assertEquals(1, results.PassCount(), t)
	assertEquals(0, results.FailCount(), t)
}

func Test__When_a_spec_contains_failing_asserts__Then_the_spec_fails(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Then(1).Should.Equal(2)
	})
	assertEquals(0, results.PassCount(), t)
	assertEquals(1, results.FailCount(), t)
}


func Test__When_a_spec_has_passing_SHOULDs__Then_its_children_are_executed(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Then(1).Should.Equal(1)
		c.Specify("Child", func() {
		})
	})
	assertEquals(2, results.TotalCount(), t)
}

func Test__When_a_spec_has_failing_SHOULDs__Then_its_children_are_executed(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Then(1).Should.Equal(2)
		c.Specify("Child", func() {
		})
	})
	assertEquals(2, results.TotalCount(), t)
}

func Test__When_a_spec_has_passing_MUSTs__Then_its_children_are_executed(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Then(1).Must.Equal(1)
		c.Specify("Child", func() {
		})
	})
	assertEquals(2, results.TotalCount(), t)
}

func Test__When_a_spec_has_failing_MUSTs__Then_its_children_are_NOT_executed(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Then(1).Must.Equal(2)
		c.Specify("Child", func() {
		})
	})
	assertEquals(1, results.TotalCount(), t)
}


// Location

func Test__The_location_of_a_failed_THEN_is_reported(t *testing.T) {
	results := runSpec(func(c Context) {
		c.Then(1).Should.Equal(2)
	})
	assertErrorIsInFile("asserts_test.go", results, t)
}

