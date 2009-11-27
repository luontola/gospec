// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"testing";
)


// Specs with children, but without siblings

func Test__Given_a_spec_with_no_children__When_it_is_run_initially__Then_the_root_is_executed(t *testing.T) {
	resetTestSpy();

	r := NewRootSpecRunner("DummySpecWithNoChildren", DummySpecWithNoChildren);
	r.runInContext(newInitialContext());

	assertTestSpyHas("root", t);
}

func Test__Given_a_spec_with_one_child__When_it_is_run_initially__Then_the_child_is_executed(t *testing.T) {
	resetTestSpy();

	r := NewRootSpecRunner("DummySpecWithOneChild", DummySpecWithOneChild);
	r.runInContext(newInitialContext());

	assertTestSpyHas("root,a", t);
}

func Test__Given_a_spec_with_nested_children__When_it_is_run_initially__Then_the_nested_children_are_executed(t *testing.T) {
	resetTestSpy();

	r := NewRootSpecRunner("DummySpecWithNestedChildren", DummySpecWithNestedChildren);
	r.runInContext(newInitialContext());

	assertTestSpyHas("root,a,aa", t);
}


func DummySpecWithNoChildren(c *Context) {
	testSpy += "root";
}

func DummySpecWithOneChild(c *Context) {
	testSpy += "root";
	c.Specify("Child A", func() {
		testSpy += ",a";
	});
}

func DummySpecWithNestedChildren(c *Context) {
	testSpy += "root";
	c.Specify("Child A", func() {
		testSpy += ",a";
		c.Specify("Child AA", func() {
			testSpy += ",aa";
		});
	});
}


// Specs with siblings, execute only one sibling at a time

func Test__Given_a_spec_with_two_children__When_it_is_run_initially__Then_the_1st_child_is_executed(t *testing.T) {
	resetTestSpy();

	r := NewRootSpecRunner("DummySpecWithTwoChildren", DummySpecWithTwoChildren);
	r.runInContext(newInitialContext());

	assertTestSpyHas("root,a", t);
}

func Test__Given_a_spec_with_two_children__When_the_1st_child_is_run_explicitly__Then_the_1st_child_is_executed(t *testing.T) {
	resetTestSpy();

	r := NewRootSpecRunner("DummySpecWithTwoChildren", DummySpecWithTwoChildren);
	r.runInContext(newExplicitContext([]int{0}));

	assertTestSpyHas("root,a", t);
}

func Test__Given_a_spec_with_two_children__When_the_2nd_child_is_run_explicitly__Then_the_2nd_child_is_executed(t *testing.T) {
	resetTestSpy();

	r := NewRootSpecRunner("DummySpecWithTwoChildren", DummySpecWithTwoChildren);
	r.runInContext(newExplicitContext([]int{1}));

	assertTestSpyHas("root,b", t);
}

func DummySpecWithTwoChildren(c *Context) {
	testSpy += "root";
	c.Specify("Child A", func() {
		testSpy += ",a";
	});
	c.Specify("Child B", func() {
		testSpy += ",b";
	});
}

