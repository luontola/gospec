// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/vector"
	"sort"
	"testing"
)


// Specs with children, but without siblings

func Test__Given_a_spec_with_no_children__When_it_is_run_initially__Then_the_root_is_executed(t *testing.T) {
	runSpec("DummySpecWithNoChildren", DummySpecWithNoChildren, newInitialContext())
	assertTestSpyHas("root", t)
}

func Test__Given_a_spec_with_one_child__When_it_is_run_initially__Then_the_child_is_executed(t *testing.T) {
	runSpec("DummySpecWithOneChild", DummySpecWithOneChild, newInitialContext())
	assertTestSpyHas("root,a", t)
}

func Test__Given_a_spec_with_nested_children__When_it_is_run_initially__Then_the_nested_children_are_executed(t *testing.T) {
	runSpec("DummySpecWithNestedChildren", DummySpecWithNestedChildren, newInitialContext())
	assertTestSpyHas("root,a,aa", t)
}


// Specs with siblings, execute only one sibling at a time

func Test__Given_a_spec_with_two_children__When_it_is_run_initially__Then_the_1st_child_is_executed(t *testing.T) {
	runSpec("DummySpecWithTwoChildren", DummySpecWithTwoChildren, newInitialContext())
	assertTestSpyHas("root,a", t)
}

func Test__Given_a_spec_with_two_children__When_the_1st_child_is_run_explicitly__Then_the_1st_child_is_executed(t *testing.T) {
	runSpec("DummySpecWithTwoChildren", DummySpecWithTwoChildren, newExplicitContext([]int{0}))
	assertTestSpyHas("root,a", t)
}

func Test__Given_a_spec_with_two_children__When_the_2nd_child_is_run_explicitly__Then_the_2nd_child_is_executed(t *testing.T) {
	runSpec("DummySpecWithTwoChildren", DummySpecWithTwoChildren, newExplicitContext([]int{1}))
	assertTestSpyHas("root,b", t)
}


// Specs with nested siblings, execute eventually all siblings, one at a time

func Test__Given_a_spec_with_multiple_nested_children__When_it_is_run_fully__Then_all_the_children_are_executed_in_isolation(t *testing.T) {
	r := NewRunner()
	r.AddSpec("DummySpecWithMultipleNestedChildren", DummySpecWithMultipleNestedChildren)

	// Execute manually instead of calling Run(), in order to avoid running
	// the specs multi-threadedly, which would mess up the test spy.
	runs := new(vector.StringVector)
	for r.hasScheduledTasks() {
		resetTestSpy()
		r.executeNextScheduledTask()
		runs.Push(testSpy)
	}
	sort.Sort(runs)

	assertEquals(5, runs.Len(), t)
	assertEquals("root,a,aa", runs.At(0), t)
	assertEquals("root,a,ab", runs.At(1), t)
	assertEquals("root,b,ba", runs.At(2), t)
	assertEquals("root,b,bb", runs.At(3), t)
	assertEquals("root,b,bc", runs.At(4), t)
}

