// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"testing";
)


// When specs are run, they should report which specs were executed now and which were postponed

func Test__Executed_specs_are_reported(t *testing.T) {
	result := runSpec("DummySpecWithTwoChildren", DummySpecWithTwoChildren, newInitialContext());
	
	executed := result.executedSpecs;
	assertEquals(2, len(executed), t);
	assertEquals("DummySpecWithTwoChildren", executed[0].name, t);
	assertEquals("Child A", executed[1].name, t);
}

func Test__Postponed_specs_are_reported(t *testing.T) {
	result := runSpec("DummySpecWithTwoChildren", DummySpecWithTwoChildren, newInitialContext());
	
	postponed := result.postponedSpecs;
	assertEquals(1, len(postponed), t);
	assertEquals("Child B", postponed[0].name, t);
}

func Test__Previously_executed_specs_are_NOT_reported(t *testing.T) {
	result := runSpec("DummySpecWithTwoChildren", DummySpecWithTwoChildren, newExplicitContext([]int{1}));
	
	executed := result.executedSpecs;
	assertEquals(2, len(executed), t);
	assertEquals("DummySpecWithTwoChildren", executed[0].name, t);
	assertEquals("Child B", executed[1].name, t);

	postponed := result.postponedSpecs;
	assertEquals(0, len(postponed), t);
}


// Scheduling specs for execution

func Test__The_postponed_specs_are_scheduled_for_execution_until_they_all_have_been_executed(t *testing.T) {
	r := NewSpecRunner();
	r.AddSpec("DummySpecWithTwoChildren", DummySpecWithTwoChildren);
	r.Run();
	
	runCounts := countSpecNames(r.executed);
	assertEquals(3, len(runCounts), t);
	assertEquals(2, runCounts["DummySpecWithTwoChildren"], t);
	assertEquals(1, runCounts["Child A"], t);
	assertEquals(1, runCounts["Child B"], t);
}

func Test__Multiple_specs_can_be_executed_in_one_batch(t *testing.T) {
	r := NewSpecRunner();
	r.AddSpec("DummySpecWithOneChild", DummySpecWithOneChild);
	r.AddSpec("DummySpecWithTwoChildren", DummySpecWithTwoChildren);
	r.Run();
	
	runCounts := countSpecNames(r.executed);
	assertEquals(1, runCounts["DummySpecWithOneChild"], t);
	assertEquals(2, runCounts["DummySpecWithTwoChildren"], t);
}

