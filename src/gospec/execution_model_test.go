// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"testing";
)


// Specs with children, but without siblings

func Test__Given_a_spec_with_no_children__When_it_is_run_initially__Then_the_root_is_executed(t *testing.T) {
	runSpec("DummySpecWithNoChildren", DummySpecWithNoChildren, newInitialContext());
	assertTestSpyHas("root", t);
}

func Test__Given_a_spec_with_one_child__When_it_is_run_initially__Then_the_child_is_executed(t *testing.T) {
	runSpec("DummySpecWithOneChild", DummySpecWithOneChild, newInitialContext());
	assertTestSpyHas("root,a", t);
}

func Test__Given_a_spec_with_nested_children__When_it_is_run_initially__Then_the_nested_children_are_executed(t *testing.T) {
	runSpec("DummySpecWithNestedChildren", DummySpecWithNestedChildren, newInitialContext());
	assertTestSpyHas("root,a,aa", t);
}


// Specs with siblings, execute only one sibling at a time

func Test__Given_a_spec_with_two_children__When_it_is_run_initially__Then_the_1st_child_is_executed(t *testing.T) {
	runSpec("DummySpecWithTwoChildren", DummySpecWithTwoChildren, newInitialContext());
	assertTestSpyHas("root,a", t);
}

func Test__Given_a_spec_with_two_children__When_the_1st_child_is_run_explicitly__Then_the_1st_child_is_executed(t *testing.T) {
	runSpec("DummySpecWithTwoChildren", DummySpecWithTwoChildren, newExplicitContext([]int{0}));
	assertTestSpyHas("root,a", t);
}

func Test__Given_a_spec_with_two_children__When_the_2nd_child_is_run_explicitly__Then_the_2nd_child_is_executed(t *testing.T) {
	runSpec("DummySpecWithTwoChildren", DummySpecWithTwoChildren, newExplicitContext([]int{1}));
	assertTestSpyHas("root,b", t);
}


// Specs with nested siblings, execute eventually all siblings, one at a time

func Test__Given_a_spec_with_multiple_nested_children__When_it_is_run_fully__Then_all_the_children_are_executed_in_isolation(t *testing.T) {
	runSpec("DummySpecWithMultipleNestedChildren", DummySpecWithMultipleNestedChildren, newInitialContext());
	assertTestSpyHas("root,a,aa", t);
	// TODO: replace explicit target paths with ones reported by the spec runner
	runSpec("DummySpecWithMultipleNestedChildren", DummySpecWithMultipleNestedChildren, newExplicitContext([]int{0, 1}));
	assertTestSpyHas("root,a,ab", t);
	runSpec("DummySpecWithMultipleNestedChildren", DummySpecWithMultipleNestedChildren, newExplicitContext([]int{1}));
	assertTestSpyHas("root,b,ba", t);
	runSpec("DummySpecWithMultipleNestedChildren", DummySpecWithMultipleNestedChildren, newExplicitContext([]int{1, 1}));
	assertTestSpyHas("root,b,bb", t);
	runSpec("DummySpecWithMultipleNestedChildren", DummySpecWithMultipleNestedChildren, newExplicitContext([]int{1, 2}));
	assertTestSpyHas("root,b,bc", t);
}

