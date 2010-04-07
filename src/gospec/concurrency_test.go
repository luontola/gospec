// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"math"
	"nanospec"
	"time"
)


// Executing specs using multiple threads

const (
	MILLISECOND = 1000000
	DELAY       = 50 * MILLISECOND
)

func ConcurrencySpec(c nanospec.Context) {
	r := NewRunner()
	r.AddSpec(VerySlowDummySpec)

	start := time.Nanoseconds()
	r.Run()
	end := time.Nanoseconds()
	totalTime := end - start

	// If the spec is executed single-threadedly, then it would take
	// at least 4*DELAY to execute. If executed multi-threadedly, it
	// would take at least 2*DELAY to execute, because the first spec
	// needs to be executed fully before the other specs are found, but
	// after that the other specs can be executed in parallel.
	expectedMaxTime := int64(math.Floor(2.9 * DELAY))

	if totalTime > expectedMaxTime {
		c.Errorf("Expected the run to take less than %v ms but it took %v ms",
			expectedMaxTime/MILLISECOND, totalTime/MILLISECOND)
	}

	runCounts := countSpecNames(r.executed)
	c.Expect(runCounts["Child A"]).Equals(1)
	c.Expect(runCounts["Child B"]).Equals(1)
	c.Expect(runCounts["Child C"]).Equals(1)
	c.Expect(runCounts["Child D"]).Equals(1)
}

func VerySlowDummySpec(c Context) {
	c.Specify("A very slow test setup", func() {
		time.Sleep(DELAY)
		c.Specify("Child A", func() {
		})
		c.Specify("Child B", func() {
		})
		c.Specify("Child C", func() {
		})
		c.Specify("Child D", func() {
		})
	})
}
