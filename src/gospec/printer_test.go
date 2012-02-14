// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"bytes"
	"github.com/orfjackal/nanospec.go/src/nanospec"
	"strings"
)

var noErrors = []*Error{}
var someError = []*Error{newError(OtherError, "some error", "", []*Location{})}

func PrinterSpec(c nanospec.Context) {
	trim := strings.TrimSpace
	out := new(bytes.Buffer)
	p := NewPrinter(SimplePrintFormat(out))

	c.Specify("When showing the summary", func() {
		p.ShowAll()
		p.ShowSummary()

		c.Specify("then the summary is printed", func() {
			p.VisitSpec(0, "Passing 1", noErrors)
			p.VisitSpec(0, "Passing 2", noErrors)
			p.VisitSpec(0, "Failing", someError)
			p.VisitEnd(2, 1)
			c.Expect(trim(out.String())).Equals(trim(`
- Passing 1
- Passing 2
- Failing [FAIL]
*** some error

3 specs, 1 failures
`))
		})
	})
	c.Specify("When hiding the summary", func() {
		p.ShowAll()
		p.HideSummary()

		c.Specify("then the summary is not printed", func() {
			p.VisitSpec(0, "Passing 1", noErrors)
			p.VisitSpec(0, "Passing 2", noErrors)
			p.VisitSpec(0, "Failing", someError)
			p.VisitEnd(2, 1)
			c.Expect(trim(out.String())).Equals(trim(`
- Passing 1
- Passing 2
- Failing [FAIL]
*** some error
`))
		})
	})

	c.Specify("When showing all specs", func() {
		p.ShowAll()

		c.Specify("then passing and failing specs are printed", func() {
			p.VisitSpec(0, "Passing", noErrors)
			p.VisitSpec(0, "Failing", someError)
			c.Expect(trim(out.String())).Equals(trim(`
- Passing
- Failing [FAIL]
*** some error
`))
		})
	})
	c.Specify("When showing only failing specs", func() {
		p.ShowOnlyFailing()

		c.Specify("then only failing specs are printed", func() {
			p.VisitSpec(0, "Passing", noErrors)
			p.VisitSpec(0, "Failing", someError)
			c.Expect(trim(out.String())).Equals(trim(`
- Failing [FAIL]
*** some error
`))
		})

		c.Specify("then the parents of failing specs are printed", func() {
			p.VisitSpec(0, "Passing parent", noErrors)
			p.VisitSpec(1, "Failing child", someError)
			c.Expect(trim(out.String())).Equals(trim(`
- Passing parent
  - Failing child [FAIL]
*** some error
`))
		})

		c.Specify("Case: passing parent with many failing children; should print the parent only once", func() {
			p.VisitSpec(0, "Passing parent", noErrors)
			p.VisitSpec(1, "Failing child A", someError)
			p.VisitSpec(1, "Failing child B", someError)
			c.Expect(trim(out.String())).Equals(trim(`
- Passing parent
  - Failing child A [FAIL]
*** some error
  - Failing child B [FAIL]
*** some error
`))
		})

		c.Specify("Case: failing parent with a failing grandchild; should print the child in the middle", func() {
			p.VisitSpec(0, "Failing parent", someError)
			p.VisitSpec(1, "Passing child", noErrors)
			p.VisitSpec(2, "Failing grandchild", someError)
			c.Expect(trim(out.String())).Equals(trim(`
- Failing parent [FAIL]
*** some error
  - Passing child
    - Failing grandchild [FAIL]
*** some error
`))
		})

		c.Specify("Case: failing parent and ghosts of unrelated specs; should not print unrelated specs", func() {
			p.VisitSpec(0, "Don't show me 0", noErrors)
			p.VisitSpec(1, "Don't show me 1", noErrors)
			p.VisitSpec(2, "Don't show me 2", noErrors)
			p.VisitSpec(0, "Failing parent", someError)
			p.VisitSpec(1, "Failing child", someError)
			c.Expect(trim(out.String())).Equals(trim(`
- Failing parent [FAIL]
*** some error
  - Failing child [FAIL]
*** some error
`))
		})
	})
}
