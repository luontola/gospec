// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"nanospec"
)


func ExpectationsSpec(c nanospec.Context) {

	c.Specify("When a spec has passing expectations or assumptions", func() {
		results := runSpec(func(c Context) {
			c.Expect(1, Equals, 1)
			c.Assume(1, Equals, 1)
			c.Specify("Child", func() {})
		})

		c.Specify("then the spec passes", func() {
			c.Expect(results.FailCount()).Equals(0)
		})
		c.Specify("then its children are executed", func() {
			c.Expect(results.TotalCount()).Equals(2)
		})
	})

	c.Specify("When a spec has failing expectations", func() {
		results := runSpec(func(c Context) {
			c.Expect(1, Equals, 2)
			c.Specify("Child", func() {})
		})

		c.Specify("then the spec fails", func() {
			c.Expect(results.FailCount()).Equals(1)
		})
		c.Specify("then its children are executed", func() {
			c.Expect(results.TotalCount()).Equals(2)
		})
	})

	c.Specify("When a spec has failing assumptions", func() {
		results := runSpec(func(c Context) {
			c.Assume(1, Equals, 2)
			c.Specify("Child", func() {})
		})

		c.Specify("then the spec fails", func() {
			c.Expect(results.FailCount()).Equals(1)
		})
		c.Specify("then its children are NOT executed", func() {
			c.Expect(results.TotalCount()).Equals(1)
		})
	})

	c.Specify("The location of a failed expectation is reported", func() {
		results := runSpec(func(c Context) {
			c.Expect(1, Equals, 2)
		})
		c.Expect(fileOfError(results)).Equals("expectations_test.go")
	})
	c.Specify("The location of a failed assumption is reported", func() {
		results := runSpec(func(c Context) {
			c.Assume(1, Equals, 2)
		})
		c.Expect(fileOfError(results)).Equals("expectations_test.go")
	})
}

func fileOfError(results *ResultCollector) string {
	file := ""
	for spec := range results.sortedRoots() {
		error := spec.errors.Front().Value.(*Error)
		file = error.Location.File
	}
	return file
}
