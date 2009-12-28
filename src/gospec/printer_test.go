// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"bytes"
	"testing"
)

var noErrors = []*Error{}
var someError = []*Error{newError("some error", currentLocation())}


// Showing the summary

func Test__When_printing_a_summary__Then_a_summary_is_printed(t *testing.T) {
	out := new(bytes.Buffer)
	
	p := newPrinter(out)
	p.ShowAll()
	p.ShowSummary()
	p.VisitSpec(0, "Passing 1", noErrors)
	p.VisitSpec(0, "Passing 2", noErrors)
	p.VisitSpec(0, "Failing", someError)
	p.VisitEnd(2, 1)
	
	assertEqualsTrim(`
- Passing 1
- Passing 2
- Failing [FAIL]
    some error

3 specs, 1 failures
	`, out.String(), t)
}

func Test__When_not_printing_a_summary__Then_a_summary_is_not_printed(t *testing.T) {
	out := new(bytes.Buffer)
	
	p := newPrinter(out)
	p.ShowAll()
	p.HideSummary()
	p.VisitSpec(0, "Passing 1", noErrors)
	p.VisitSpec(0, "Passing 2", noErrors)
	p.VisitSpec(0, "Failing", someError)
	p.VisitEnd(2, 1)
	
	assertEqualsTrim(`
- Passing 1
- Passing 2
- Failing [FAIL]
    some error
	`, out.String(), t)
}


// Showing only failures

func Test__When_printing_all_specs__Then_passing_and_failing_are_printed(t *testing.T) {
	out := new(bytes.Buffer)
	
	p := newPrinter(out)
	p.ShowAll()
	p.VisitSpec(0, "Passing", noErrors)
	p.VisitSpec(0, "Failing", someError)
	
	assertEqualsTrim(`
- Passing
- Failing [FAIL]
    some error
	`, out.String(), t)
}

func Test__When_printing_only_failing_specs__Then_only_failing_are_printed(t *testing.T) {
	out := new(bytes.Buffer)
	
	p := newPrinter(out)
	p.ShowOnlyFailing()
	p.VisitSpec(0, "Passing", noErrors)
	p.VisitSpec(0, "Failing", someError)
	
	assertEqualsTrim(`
- Failing [FAIL]
    some error
	`, out.String(), t)
}

func Test__When_printing_only_failing_specs__Then_the_parents_of_failing_specs_are_printed(t *testing.T) {
	out := new(bytes.Buffer)
	
	p := newPrinter(out)
	p.ShowOnlyFailing()
	p.VisitSpec(0, "Passing parent", noErrors)
	p.VisitSpec(1, "Failing child", someError)
	
	assertEqualsTrim(`
- Passing parent
  - Failing child [FAIL]
      some error
	`, out.String(), t)
}

func Test__Passing_parent_with_many_failing_children(t *testing.T) {
	out := new(bytes.Buffer)
	
	p := newPrinter(out)
	p.ShowOnlyFailing()
	p.VisitSpec(0, "Passing parent", noErrors)
	p.VisitSpec(1, "Failing child A", someError)
	p.VisitSpec(1, "Failing child B", someError)
	
	assertEqualsTrim(`
- Passing parent
  - Failing child A [FAIL]
      some error
  - Failing child B [FAIL]
      some error
	`, out.String(), t)
}

func Test__Failing_parent_with_a_failing_grandchild(t *testing.T) {
	out := new(bytes.Buffer)
	
	p := newPrinter(out)
	p.ShowOnlyFailing()
	p.VisitSpec(0, "Failing parent", someError)
	p.VisitSpec(1, "Passing child", noErrors)
	p.VisitSpec(2, "Failing grandchild", someError)
	
	assertEqualsTrim(`
- Failing parent [FAIL]
    some error
  - Passing child
    - Failing grandchild [FAIL]
        some error
	`, out.String(), t)
}

func Test__Failing_parent_and_ghosts_of_unrelated_specs(t *testing.T) {
	out := new(bytes.Buffer)
	
	p := newPrinter(out)
	p.ShowOnlyFailing()
	p.VisitSpec(0, "Don't show me 0", noErrors)
	p.VisitSpec(1, "Don't show me 1", noErrors)
	p.VisitSpec(2, "Don't show me 2", noErrors)
	p.VisitSpec(0, "Failing parent", someError)
	p.VisitSpec(1, "Failing child", someError)
	
	assertEqualsTrim(`
- Failing parent [FAIL]
    some error
  - Failing child [FAIL]
      some error
	`, out.String(), t)
}

