// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"testing";
	"fmt";
)




func Test__When_no_specs_are_executed__Then_the_report_has_in_total_zero_specs(t *testing.T) {
	report := newResultCollector();
	
	assertEquals(0, report.TotalCount(), t);
}

func Test__When_one_root_spec_is_executed__Then_the_report_has_in_total_one_spec(t *testing.T) {
	report := newResultCollector();
	report.Update(newSpecRun("RootSpec", nil, nil));
	
	assertEquals(1, report.TotalCount(), t);
}

func Test__When_one_root_spec_is_executed__Then_the_report_has_one_root_spec_with_no_children(t *testing.T) {
	// TODO: refactor so that the tree structure is visible
	report := newResultCollector();
	report.Update(newSpecRun("RootSpec", nil, nil));
	
	roots := report.Roots();
	r1 := <-roots;
	assertHasNoMore(roots, t);
	
	assertSpecHasName("RootSpec", r1, t);
	assertHasNoMore(r1.Children(), t);
}

func Test__When_many_root_specs_are_executed__Then_the_report_has_many_root_specs(t *testing.T) {
	// TODO: refactor so that the tree structure is visible
	report := newResultCollector();
	report.Update(newSpecRun("RootSpec1", nil, nil));
	report.Update(newSpecRun("RootSpec2", nil, nil));
	
	roots := report.Roots();
	r1 := <-roots;
	r2 := <-roots;
	assertHasNoMore(roots, t);

	assertSpecHasName("RootSpec1", r1, t);
	assertSpecHasName("RootSpec2", r2, t);
}

func Test__When_nested_specs_are_executed__Then_the_root_spec_has_children(t *testing.T) {
	// TODO: refactor so that the tree structure is visible
	report := newResultCollector();
	s1 := newSpecRun("RootSpec", nil, nil);
	s2 := newSpecRun("Child", nil, s1);
	report.Update(s1);
	report.Update(s2);
	
	roots := report.Roots();
	r1 := <-roots;
	assertHasNoMore(roots, t);

	children := r1.Children();
	r1c1 := <-children;
	assertHasNoMore(children, t);

	assertSpecHasName("RootSpec", r1, t);
	assertSpecHasName("Child", r1c1, t);
}




func Test__The_total_number_of_specs_is_reported(t *testing.T) {
	/*
	r := NewSpecRunner();
	r.AddSpec("DummySpecWithOneFailure", DummySpecWithOneFailure);
	r.Run();
	assertEquals(3, r.TotalCount(), t);
	*/
}

func DummySpecWithOneFailure(c *Context) {
	c.Specify("Failing spec", func() {
	});
	c.Specify("Passing spec", func() {
	});
}

func assertHasNoMore(iter <-chan *specResult, t *testing.T) {
	assertEquals(true, nil == <-iter, t);
}

func assertSpecHasName(name string, spec *specResult, t *testing.T) {
	if spec == nil {
		t.Error(fmt.Sprintf("Expected a spec with name '%v' but the spec was nil", name));
	} else {
		assertEquals(name, spec.name, t);
	}
}

