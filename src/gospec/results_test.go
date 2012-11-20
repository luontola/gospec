// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"bytes"
	"errors"
	"github.com/orfjackal/nanospec.go/src/nanospec"
	"strings"
)

func ResultsSpec(c nanospec.Context) {
	results := newResultCollector()

	c.Specify("When results have many root specs", func() {
		results.Update(newSpecRun("RootSpec2", nil, nil, nil))
		results.Update(newSpecRun("RootSpec1", nil, nil, nil))
		results.Update(newSpecRun("RootSpec3", nil, nil, nil))

		c.Specify("then the roots are sorted alphabetically", func() {
			c.Expect(results).Matches(ReportIs(`
- RootSpec1
- RootSpec2
- RootSpec3

3 specs, 0 failures
`))
		})
	})

	c.Specify("When results have many child specs", func() {
		// In tests, when a spec has many children, make sure
		// to pass a common parent instance to all the siblings.
		// Otherwise the parent's numberOfChildren is not
		// incremented and the children's paths will be wrong.

		// use names which would not sort alphabetically
		root := newSpecRun("RootSpec", nil, nil, nil)
		child1 := newSpecRun("one", nil, root, nil)
		child2 := newSpecRun("two", nil, root, nil)
		child3 := newSpecRun("three", nil, root, nil)

		// register in random order
		results.Update(root)
		results.Update(child1)

		results.Update(root)
		results.Update(child3)

		results.Update(root)
		results.Update(child2)

		c.Specify("then the children are sorted by their declaration order", func() {
			c.Expect(results).Matches(ReportIs(`
- RootSpec
  - one
  - two
  - three

4 specs, 0 failures
`))
		})
	})

	c.Specify("Case: zero specs", func() {
		c.Expect(results).Matches(ReportIs(`
0 specs, 0 failures
`))
	})
	c.Specify("Case: spec with no children", func() {
		a1 := newSpecRun("RootSpec", nil, nil, nil)
		results.Update(a1)
		c.Expect(results).Matches(ReportIs(`
- RootSpec

1 specs, 0 failures
`))
	})
	c.Specify("Case: spec with a child", func() {
		a1 := newSpecRun("RootSpec", nil, nil, nil)
		a2 := newSpecRun("Child A", nil, a1, nil)
		results.Update(a1)
		results.Update(a2)
		c.Expect(results).Matches(ReportIs(`
- RootSpec
  - Child A

2 specs, 0 failures
`))
	})
	c.Specify("Case: spec with nested children", func() {
		a1 := newSpecRun("RootSpec", nil, nil, nil)
		a2 := newSpecRun("Child A", nil, a1, nil)
		a3 := newSpecRun("Child AA", nil, a2, nil)
		results.Update(a1)
		results.Update(a2)
		results.Update(a3)
		c.Expect(results).Matches(ReportIs(`
- RootSpec
  - Child A
    - Child AA

3 specs, 0 failures
`))
	})
	c.Specify("Case: spec with multiple nested children", func() {
		runner := NewParallelRunner()
		runner.AddSpec(DummySpecWithMultipleNestedChildren)
		runner.Run()
		c.Expect(runner.Results()).Matches(ReportIs(`
- gospec.DummySpecWithMultipleNestedChildren
  - Child A
    - Child AA
    - Child AB
  - Child B
    - Child BA
    - Child BB
    - Child BC

8 specs, 0 failures
`))
	})

	c.Specify("When specs fail", func() {
		a1 := newSpecRun("Failing", nil, nil, nil)
		a1.AddError(newError(OtherError, "X did not equal Y", "", []*Location{}))
		results.Update(a1)

		b1 := newSpecRun("Passing", nil, nil, nil)
		b2 := newSpecRun("Child failing", nil, b1, nil)
		b2.AddError(newError(OtherError, "moon was not cheese", "", []*Location{}))
		results.Update(b1)
		results.Update(b2)

		c.Specify("then the errors are reported", func() {
			c.Expect(results).Matches(ReportIs(`
- Failing [FAIL]
*** X did not equal Y
- Passing
  - Child failing [FAIL]
*** moon was not cheese

3 specs, 2 failures
`))
		})
	})
	c.Specify("When spec passes on 1st run but fails on 2nd run", func() {
		i := 0
		runner := NewParallelRunner()
		runner.AddNamedSpec("RootSpec", func(c Context) {
			if i == 1 {
				c.Expect(10, Equals, 20)
			}
			i++
			c.Specify("Child A", func() {})
			c.Specify("Child B", func() {})
		})
		runner.Run()

		c.Specify("then the error is reported", func() {
			c.Expect(runner.Results()).Matches(ReportIs(`
- RootSpec [FAIL]
*** Expected: equals “20”
         got: “10”
    at results_test.go
  - Child A
  - Child B

3 specs, 1 failures
`))
		})
	})
	c.Specify("When root spec fails sporadically", func() {
		runner := NewParallelRunner()
		runner.AddNamedSpec("RootSpec", func(c Context) {
			i := 0
			c.Specify("Child A", func() {
				i = 1
			})
			c.Specify("Child B", func() {
				i = 2
			})
			c.Expect(10, Equals, 20)   // stays same - will be reported once
			c.Expect(10+i, Equals, 20) // changes - will be reported many times
		})
		runner.Run()

		c.Specify("then the errors are merged together", func() {
			c.Expect(runner.Results()).Matches(ReportIs(`
- RootSpec [FAIL]
*** Expected: equals “20”
         got: “10”
    at results_test.go
*** Expected: equals “20”
         got: “11”
    at results_test.go
*** Expected: equals “20”
         got: “12”
    at results_test.go
  - Child A
  - Child B

3 specs, 1 failures
`))
		})
	})
	c.Specify("When non-root spec fails sporadically", func() {
		runner := NewParallelRunner()
		runner.AddNamedSpec("RootSpec", func(c Context) {
			c.Specify("Failing", func() {
				i := 0
				c.Specify("Child A", func() {
					i = 1
				})
				c.Specify("Child B", func() {
					i = 2
				})
				c.Expect(10, Equals, 20)   // stays same - will be reported once
				c.Expect(10+i, Equals, 20) // changes - will be reported many times
			})
		})
		runner.Run()

		c.Specify("then the errors are merged together", func() {
			c.Expect(runner.Results()).Matches(ReportIs(`
- RootSpec
  - Failing [FAIL]
*** Expected: equals “20”
         got: “10”
    at results_test.go
*** Expected: equals “20”
         got: “11”
    at results_test.go
*** Expected: equals “20”
         got: “12”
    at results_test.go
    - Child A
    - Child B

4 specs, 1 failures
`))
		})
	})

	c.Specify("When an expectation gives an error", func() {
		runner := NewParallelRunner()
		runner.AddNamedSpec("RootSpec", func(c Context) {
			c.Expect(1, IsWithin(0.1), 1.0)
		})
		runner.Run()

		c.Specify("the error is reported as-is", func() {
			c.Expect(runner.Results()).Matches(ReportIs(`
- RootSpec [FAIL]
*** type error: expected a float, but was “1” of type “int”
    at results_test.go

1 specs, 1 failures
`))
		})
	})

	c.Specify("When a spec panics", func() {
		runner := NewParallelRunner()
		runner.AddNamedSpec("RootSpec", func(c Context) {
			c.Specify("Child A", func() {
				boom2()
			})
		})
		runner.Run()

		c.Specify("then the panic's stack trace is reported", func() {
			c.Expect(runner.Results()).Matches(ReportIs(`
- RootSpec
  - Child A [FAIL]
*** panic: boom!
    at recover_test.go
    at recover_test.go
    at recover_test.go
    at results_test.go

2 specs, 1 failures
`))
		})
	})

	c.Specify("When a root spec panics", func() {
		runner := NewParallelRunner()
		runner.AddNamedSpec("RootSpec", func(c Context) {
			boom2()
		})

		runner.Run()
		c.Specify("the bootstrap code in runner.go does not show up in the stack trace", func() {
			c.Expect(runner.Results()).Matches(ReportIs(`
- RootSpec [FAIL]
*** panic: boom!
    at recover_test.go
    at recover_test.go
    at recover_test.go
    at results_test.go

1 specs, 1 failures
`))
		})
	})
}

func ReportIs(expected string) nanospec.Matcher {
	return func(v interface{}) error {
		actual := strings.TrimSpace(resultToString(v.(*ResultCollector)))
		expected = strings.TrimSpace(expected)
		if actual != expected {
			return errors.New("Expected report:\n" + expected + "\n\nBut was:\n" + actual)
		}
		return nil
	}
}

func ReportContains(needle string) nanospec.Matcher {
	return func(v interface{}) error {
		actual := resultToString(v.(*ResultCollector))
		found := strings.Index(actual, needle) >= 0
		if !found {
			return errors.New("Expected report to contain:\n" + needle + "\n\nBut report was:\n" + actual)
		}
		return nil
	}
}

// TODO: convert all the tests to use a matcher like this, which takes itself care of running the spec
func SpecsReportContains(needle string) nanospec.Matcher {
	return func(v interface{}) error {
		spec := v.(func(Context))

		runner := NewParallelRunner()
		runner.AddNamedSpec("RootSpec", spec)
		runner.Run()

		return ReportContains(needle)(runner.Results())
	}
}

func resultToString(result *ResultCollector) string {
	out := new(bytes.Buffer)
	result.Visit(NewPrinter(SimplePrintFormat(out)))
	return out.String()
}
