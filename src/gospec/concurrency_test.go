// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"testing";
	"time";
)


// Executing specs using multiple threads

const (
	millisecond = 1000000;
	delay = 50 * millisecond;
)

func Test__Specs_are_executed_concurrently_on_multiple_threads(t *testing.T) {
	r := NewSpecRunner();
	r.AddSpec("VerySlowDummySpec", VerySlowDummySpec);
	
	start := time.Nanoseconds();
	r.Run();
	end := time.Nanoseconds();
	totalTime := end - start;
	
	// If the spec is executed single-threadedly, then it would take
	// at least "4 * delay" to execute. If executed multi-threadedly, it
	// would take at least "2 * delay" to execute, because the first spec
	// needs to be executed fully before the other specs are found, but
	// after that the other specs can be executed in parallel.
	expectedMaxTime := int64(2.5 * delay);
	
	if totalTime > expectedMaxTime {
		t.Errorf("Expected to take less than %v ms but it took %v ms", 
			expectedMaxTime / millisecond,
			totalTime / millisecond);
	}
	
	runCounts := countSpecNames(r.executed);
	assertEquals(1, runCounts["Child A"], t);
	assertEquals(1, runCounts["Child B"], t);
	assertEquals(1, runCounts["Child C"], t);
	assertEquals(1, runCounts["Child D"], t);
}

func VerySlowDummySpec(c *Context) {
	c.Specify("A very slow test setup", func() {
		time.Sleep(delay);
		c.Specify("Child A", func() {
		});
		c.Specify("Child B", func() {
		});
		c.Specify("Child C", func() {
		});
		c.Specify("Child D", func() {
		});
	});
}

