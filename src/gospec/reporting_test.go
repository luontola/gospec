// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"testing"
)


func Test__When_results_have_many_root_specs__Then_they_are_sorted_alphabetically(t *testing.T) {
	results := newResultCollector()

	// register in reverse order
	a1 := newSpecRun("RootSpec2", nil, nil)
	results.Update(a1)

	b2 := newSpecRun("RootSpec1", nil, nil)
	results.Update(b2)

	// expect roots to be in alphabetical order
	assertReportIs(results, `
- RootSpec1
- RootSpec2
	`, 2, 0, t)
}

func Test__When_results_have_many_child_specs__Then_they_are_sorted_by_their_declaration_order(t *testing.T) {
	results := newResultCollector()

	// In tests, when a spec has many children, make sure
	// to pass a common parent instance to all the siblings.
	// Otherwise the parent's numberOfChildren is not
	// incremented and the children's paths will be wrong.

	// use names which would not sort alphabetically
	root := newSpecRun("RootSpec", nil, nil)
	child1 := newSpecRun("one", nil, root)
	child2 := newSpecRun("two", nil, root)
	child3 := newSpecRun("three", nil, root)

	// register in random order
	results.Update(root)
	results.Update(child1)

	results.Update(root)
	results.Update(child3)

	results.Update(root)
	results.Update(child2)

	// expect children to be in declaration order
	assertReportIs(results, `
- RootSpec
  - one
  - two
  - three
	`, 4, 0, t)
}

func Test__Collecting_results_of_zero_specs(t *testing.T) {
	results := newResultCollector()

	assertReportIs(results, `
	`, 0, 0, t)
}

func Test__Collecting_results_of_a_spec_with_no_children(t *testing.T) {
	results := newResultCollector()

	a1 := newSpecRun("RootSpec", nil, nil)
	results.Update(a1)

	assertReportIs(results, `
- RootSpec
	`, 1, 0, t)
}

func Test__Collecting_results_of_a_spec_with_a_child(t *testing.T) {
	results := newResultCollector()

	a1 := newSpecRun("RootSpec", nil, nil)
	a2 := newSpecRun("Child A", nil, a1)
	results.Update(a1)
	results.Update(a2)

	assertReportIs(results, `
- RootSpec
  - Child A
	`, 2, 0, t)
}

func Test__Collecting_results_of_a_spec_with_nested_children(t *testing.T) {
	results := newResultCollector()

	a1 := newSpecRun("RootSpec", nil, nil)
	a2 := newSpecRun("Child A", nil, a1)
	a3 := newSpecRun("Child AA", nil, a2)
	results.Update(a1)
	results.Update(a2)
	results.Update(a3)

	assertReportIs(results, `
- RootSpec
  - Child A
    - Child AA
	`, 3, 0, t)
}

func Test__Collecting_results_of_a_spec_with_multiple_nested_children(t *testing.T) {
	runner := NewRunner()
	runner.AddSpec("DummySpecWithMultipleNestedChildren", DummySpecWithMultipleNestedChildren)
	runner.Run()

	assertReportIs(runner.compileResults(), `
- DummySpecWithMultipleNestedChildren
  - Child A
    - Child AA
    - Child AB
  - Child B
    - Child BA
    - Child BB
    - Child BC
	`, 8, 0, t)
}

func Test__Collecting_results_of_failing_specs(t *testing.T) {
	results := newResultCollector()

	a1 := newSpecRun("Failing", nil, nil)
	a1.AddError("X did not equal Y")
	results.Update(a1)

	b1 := newSpecRun("Passing", nil, nil)
	b2 := newSpecRun("Child failing", nil, b1)
	b2.AddError("moon was not cheese")
	results.Update(b1)
	results.Update(b2)

	assertReportIs(results, `
- Failing [FAIL]
    X did not equal Y
- Passing
  - Child failing [FAIL]
      moon was not cheese
	`, 1, 2, t)
}


func assertReportIs(results *ResultCollector, expected string, passCount int, failCount int, t *testing.T) {
	report := newReportPrinter()
	report.Visit(results)

	assertEquals(passCount, results.PassCount(), t)
	assertEquals(failCount, results.FailCount(), t)
	assertEquals(passCount+failCount, results.TotalCount(), t)
	assertEqualsTrim(expected, report.String(), t)
}

