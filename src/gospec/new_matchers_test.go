// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"testing"
	"os"
)


func Test__Placeholders_in_error_messages_are_replaced_with_variables(t *testing.T) {
	m := Errorf("%v x %v", 1, 2)
	assertEquals("1 x 2", m.String(), t)
}

func Test__Positive_expectation_failures_are_reported_with_the_positive_message(t *testing.T) {
	log := new(spyErrorLogger)
	m := newMatcherAdapter(nil, log)
	
	m.Expect(1, dummyMatcher, 1)
	log.ShouldHaveNoErrors(t)
	
	m.Expect(1, dummyMatcher, 2)
	log.ShouldHaveTheError("Positive failure: 1, 2", t)
}

func Test__Negative_expectation_failures_are_reported_with_the_negative_message(t *testing.T) {
	log := new(spyErrorLogger)
	m := newMatcherAdapter(nil, log)
	
	m.Expect(1, Not(dummyMatcher), 2)
	log.ShouldHaveNoErrors(t)
	
	m.Expect(1, Not(dummyMatcher), 1)
	log.ShouldHaveTheError("Negative failure: 1, 1", t)
}

func dummyMatcher(actual interface{}, expected interface{}) (ok bool, pos os.Error, neg os.Error) {
	ok = actual == expected
	pos = Errorf("Positive failure: %v, %v", actual, expected)
	neg = Errorf("Negative failure: %v, %v", actual, expected)
	return
}

