// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/vector"
)


func runSpec(spec func(Context)) *ResultCollector {
	r := NewRunner()
	r.AddNamedSpec("RootSpec", spec)
	r.Run()
	return r.Results()
}

func runSpecWithContext(closure func(Context), context *taskContext) *taskResult {
	resetTestSpy()
	r := NewRunner()
	return r.execute("RootSpec", closure, context)
}

func countSpecNames(specs *vector.Vector) map[string]int {
	results := make(map[string]int)
	for _, v := range *specs {
		spec := v.(*specRun)
		key := spec.name
		if _, preset := results[key]; !preset {
			results[key] = 0
		}
		results[key]++
	}
	return results
}


// Test dummies

var testSpy = ""

func resetTestSpy() {
	testSpy = ""
}

func DummySpecWithNoChildren(c Context) {
	testSpy += "root"
}

func DummySpecWithOneChild(c Context) {
	testSpy += "root"
	c.Specify("Child A", func() {
		testSpy += ",a"
	})
}

func DummySpecWithTwoChildren(c Context) {
	testSpy += "root"
	c.Specify("Child A", func() {
		testSpy += ",a"
	})
	c.Specify("Child B", func() {
		testSpy += ",b"
	})
}

func DummySpecWithNestedChildren(c Context) {
	testSpy += "root"
	c.Specify("Child A", func() {
		testSpy += ",a"
		c.Specify("Child AA", func() {
			testSpy += ",aa"
		})
	})
}

func DummySpecWithMultipleNestedChildren(c Context) {
	testSpy += "root"
	c.Specify("Child A", func() {
		testSpy += ",a"
		c.Specify("Child AA", func() {
			testSpy += ",aa"
		})
		c.Specify("Child AB", func() {
			testSpy += ",ab"
		})
	})
	c.Specify("Child B", func() {
		testSpy += ",b"
		c.Specify("Child BA", func() {
			testSpy += ",ba"
		})
		c.Specify("Child BB", func() {
			testSpy += ",bb"
		})
		c.Specify("Child BC", func() {
			testSpy += ",bc"
		})
	})
}
