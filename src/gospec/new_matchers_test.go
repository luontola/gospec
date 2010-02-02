// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"testing"
)

func Test__Placeholders_in_messages_are_replaced_with_variables(t *testing.T) {
	m := Errorf("%v x %v", 1, 2)
	assertEquals("1 x 2", m.String(), t)
}

func Test__Expect_string_EQUALS_string(t *testing.T) {
	log := new(spyErrorLogger)

	log.Expect("apple", Equals, "apple")
	log.ShouldHaveNoErrors(t)

	log.Expect("apple", Equals, "orange")
	log.ShouldHaveTheError("Expected 'orange' but was 'apple'", t)
}

func Test__Expect_should_NOT_EQUALS_string(t *testing.T) {
	log := new(spyErrorLogger)

	log.Expect("apple", Not(Equals), "orange")
	log.ShouldHaveNoErrors(t)

	log.Expect("apple", Not(Equals), "apple")
	log.ShouldHaveTheError("Did not expect 'apple' but was 'apple'", t)
}


// New matchers

func (log *spyErrorLogger) Expect(actual interface{}, matcher Matcher, expected interface{}) {
	ok, pos, _ := matcher(actual, expected)
	if !ok {
		log.AddError(newError(pos.String(), callerLocation()))
	}
}


