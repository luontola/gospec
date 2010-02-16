// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/list"
	"fmt"
	"math"
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

func Test__Errors_in_expectations_are_reported_with_the_error_message(t *testing.T) {
	log := new(spyErrorLogger)
	m := newMatcherAdapter(nil, log)
	
	m.Expect(666, dummyMatcher, 1)
	log.ShouldHaveTheError("Error: 666", t)
}

func dummyMatcher(actual interface{}, expected interface{}) (ok bool, pos os.Error, neg os.Error, err os.Error) {
	if actual.(int) == 666 {
		err = Errorf("Error: %v", actual)
		return
	}
	ok = actual == expected
	pos = Errorf("Positive failure: %v, %v", actual, expected)
	neg = Errorf("Negative failure: %v, %v", actual, expected)
	return
}


// "Equals"

func Test__Equals_matcher_on_strings(t *testing.T) {
	assertExpectation(t, "apple", Equals, "apple").Passes()
	assertExpectation(t, "apple", Equals, "orange").Fails().
		WithMessage(
			"Expected 'orange' but was 'apple'",
			"Did not expect 'orange' but was 'apple'")
}

func Test__Equals_matcher_on_ints(t *testing.T) {
	assertExpectation(t, 42, Equals, 42).Passes()
	assertExpectation(t, 42, Equals, 999).Fails()
}

func Test__Equals_matcher_on_structs(t *testing.T) {
	assertExpectation(t, DummyStruct{42, 1}, Equals, DummyStruct{42, 2}).Passes()
	assertExpectation(t, DummyStruct{42, 1}, Equals, DummyStruct{999, 2}).Fails()
}

func Test__Equals_matcher_on_struct_pointers(t *testing.T) {
	assertExpectation(t, &DummyStruct{42, 1}, Equals, &DummyStruct{42, 2}).Passes()
	assertExpectation(t, &DummyStruct{42, 1}, Equals, &DummyStruct{999, 2}).Fails()
}

type DummyStruct struct {
	value        int
	ignoredValue int
}

func (this DummyStruct) Equals(other interface{}) bool {
	switch that := other.(type) {
	case DummyStruct:
		return this.equals(&that)
	case *DummyStruct:
		return this.equals(that)
	}
	return false
}

func (this *DummyStruct) equals(that *DummyStruct) bool {
	return this.value == that.value
}

func (this DummyStruct) String() string {
	return fmt.Sprintf("DummyStruct%v", this.value)
}


// "IsSame"

func Test__IsSame_matcher(t *testing.T) {
	a1 := new(os.File)
	a2 := a1
	b := new(os.File)
	assertExpectation(t, a1, IsSame, a2).Passes()
	assertExpectation(t, a1, IsSame, b).Fails().
		WithMessage(
			fmt.Sprintf("Expected '%v' but was '%v'", b, a1),
			fmt.Sprintf("Did not expect '%v' but was '%v'", b, a1))
	assertExpectation(t, 1, IsSame, b).GivesError("Expected a pointer, but was '1' of type 'int'")
	assertExpectation(t, b, IsSame, 1).GivesError("Expected a pointer, but was '1' of type 'int'")
}


// "IsNil"

func Test__IsNil_matcher(t *testing.T) {
	assertExpectation(t, nil, IsNil).Passes() // interface value nil
	assertExpectation(t, (*int)(nil), IsNil).Passes() // typed pointer nil inside an interface value
	assertExpectation(t, new(int), IsNil).Fails()
	assertExpectation(t, 1, IsNil).Fails().
		WithMessage(
			"Expected <nil> but was '1'",
			"Did not expect <nil> but was '1'")
}


// "IsTrue"

func Test__IsTrue_matcher(t *testing.T) {
	assertExpectation(t, true, IsTrue).Passes()
	assertExpectation(t, false, IsTrue).Fails().
		WithMessage(
			"Expected 'true' but was 'false'",
			"Did not expect 'true' but was 'false'")
}


// "IsFalse"

func Test__IsFalse_matcher(t *testing.T) {
	assertExpectation(t, false, IsFalse).Passes()
	assertExpectation(t, true, IsFalse).Fails().
		WithMessage(
			"Expected 'false' but was 'true'",
			"Did not expect 'false' but was 'true'")
}


// "Satisfy"

func Test__Satisfy_matcher(t *testing.T) {
	value := 42
	assertExpectation(t, value, Satisfies, value < 100).Passes()
	assertExpectation(t, value, Satisfies, value > 100).Fails().
		WithMessage(
			"Criteria not satisfied by '42'",
			"Criteria not satisfied by '42'")
}


// "IsWithin"

func Test__IsWithin_matcher(t *testing.T) {
	value := float64(3.141)
	pi := float64(math.Pi)
	assertExpectation(t, value, IsWithin(0.001), pi).Passes()
	assertExpectation(t, value, IsWithin(0.0001), pi).Fails().
		WithMessage(
			"Expected '3.141592653589793' ± 0.0001 but was '3.141'",
			"Did not expect '3.141592653589793' ± 0.0001 but was '3.141'")
}

func Test__IsWithin_matcher_cannot_compare_ints(t *testing.T) {
	value := int(3)
	pi := float64(math.Pi)
	assertExpectation(t, value, IsWithin(0.001), pi).GivesError("Expected a float, but was '3' of type 'int'")
	assertExpectation(t, pi, IsWithin(0.001), value).GivesError("Expected a float, but was '3' of type 'int'")
}


// "Contains"

func Test__Contains_matcher(t *testing.T) {
	values := []string{"one", "two", "three"}
	assertExpectation(t, values, Contains, "one").Passes()
	assertExpectation(t, values, Contains, "two").Passes()
	assertExpectation(t, values, Contains, "three").Passes()
	assertExpectation(t, values, Contains, "four").Fails().
		WithMessage(
			"Expected 'four' to be in '[one two three]' but it was not",
			"Did not expect 'four' to be in '[one two three]' but it was")
}

func Test__Convert_array_to_array(t *testing.T) {
	values := [...]string{"one", "two", "three"}
	
	result, _ := toArray(values)
	
	assertEquals(3, len(result), t)
	assertEquals("one", result[0], t)
	assertEquals("two", result[1], t)
	assertEquals("three", result[2], t)
}

func Test__Convert_channel_to_array(t *testing.T) {
	values := list.New()
	values.PushBack("one")
	values.PushBack("two")
	values.PushBack("three")
	
	result, _ := toArray(values.Iter())
	
	assertEquals(3, len(result), t)
	assertEquals("one", result[0], t)
	assertEquals("two", result[1], t)
	assertEquals("three", result[2], t)
}

func Test__Convert_iterable_to_array(t *testing.T) {
	values := list.New()
	values.PushBack("one")
	values.PushBack("two")
	values.PushBack("three")
	
	result, _ := toArray(values)
	
	assertEquals(3, len(result), t)
	assertEquals("one", result[0], t)
	assertEquals("two", result[1], t)
	assertEquals("three", result[2], t)
}

func Test__Convert_unsupported_value_to_array(t *testing.T) {
	_, err := toArray("foo")
	
	assertEquals("Unknown type 'string', not iterable: foo", err.String(), t)
}


// "ContainsAll"

func Test__ContainsAll_matcher(t *testing.T) {
	values := []string{"one", "two", "three"}
	
	assertExpectation(t, values, ContainsAll, Values()).Passes()
	assertExpectation(t, values, ContainsAll, Values("one")).Passes()
	assertExpectation(t, values, ContainsAll, Values("three", "two")).Passes()
	assertExpectation(t, values, ContainsAll, Values("one", "two", "three")).Passes()
	
	assertExpectation(t, values, ContainsAll, Values("four")).Fails()
	assertExpectation(t, values, ContainsAll, Values("one", "four")).Fails().
		WithMessage(
			"Expected all of '[one four]' to be in '[one two three]' but they were not",
			"Did not expect all of '[one four]' to be in '[one two three]' but they were")
}


// "ContainsAny"

func Test__ContainsAny_matcher(t *testing.T) {
	values := []string{"one", "two", "three"}
	
	assertExpectation(t, values, ContainsAny, Values("one")).Passes()
	assertExpectation(t, values, ContainsAny, Values("three", "two")).Passes()
	assertExpectation(t, values, ContainsAny, Values("four", "one", "five")).Passes()
	
	assertExpectation(t, values, ContainsAny, Values()).Fails()
	assertExpectation(t, values, ContainsAny, Values("four")).Fails()
	assertExpectation(t, values, ContainsAny, Values("four", "five")).Fails().
		WithMessage(
			"Expected any of '[four five]' to be in '[one two three]' but they were not",
			"Did not expect any of '[four five]' to be in '[one two three]' but they were")
}


// "ContainsExactly"

func Test__ContainsExactly_matcher(t *testing.T) {
	values := []string{"one", "two", "three"}
	
	assertExpectation(t, values, ContainsExactly, Values("one", "two", "three")).Passes()
	assertExpectation(t, values, ContainsExactly, Values("three", "one", "two")).Passes()
	
	assertExpectation(t, values, ContainsExactly, Values()).Fails()
	assertExpectation(t, values, ContainsExactly, Values("four")).Fails()
	assertExpectation(t, values, ContainsExactly, Values("one", "two")).Fails()
	assertExpectation(t, values, ContainsExactly, Values("one", "two", "three", "four")).Fails().
		WithMessage(
			"Expected exactly '[one two three four]' to be in '[one two three]' but they were not",
			"Did not expect exactly '[one two three four]' to be in '[one two three]' but they were")
	
	// duplicate values are allowed
	values = []string{"a", "a", "b"}
	
	assertExpectation(t, values, ContainsExactly, Values("a", "a", "b")).Passes()
	assertExpectation(t, values, ContainsExactly, Values("a", "b", "a")).Passes()
	
	assertExpectation(t, values, ContainsExactly, Values("a", "b", "b")).Fails()
	assertExpectation(t, values, ContainsExactly, Values("a", "a", "a", "b")).Fails()
	assertExpectation(t, values, ContainsExactly, Values("a", "a", "b", "b")).Fails()
}



// Test utilities

func assertExpectation(t *testing.T, actual interface{}, matcher Matcher, expected ...interface{}) *matchAssert {
	ok, pos, neg, err := matcher.Match(actual, expected)
	return &matchAssert{ok, pos, neg, err, t}
}

type matchAssert struct {
	ok  bool
	pos os.Error
	neg os.Error
	err os.Error
	t *testing.T
}

func (this *matchAssert) Passes() *matchAssert {
	if this.err != nil {
		this.t.Error("expected to pass, but had an error: " + this.err.String())
	}
	if !this.ok {
		this.t.Error("expected to pass, but failed")
	}
	return this
}

func (this *matchAssert) Fails() *matchAssert {
	if this.err != nil {
		this.t.Error("expected to fail, but had an error: " + this.err.String())
	}
	if this.ok {
		this.t.Error("expected to fail, but passed")
	}
	return this
}

func (this *matchAssert) WithMessage(expectedPos string, expectedNeg string) *matchAssert {
	assertEquals(expectedPos, this.pos.String(), this.t)
	assertEquals(expectedNeg, this.neg.String(), this.t)
	return this
}

func (this *matchAssert) GivesError(expectedErr string) *matchAssert {
	if this.err == nil {
		this.t.Error("expected have an error, but did not")
	} else {
		assertEquals(expectedErr, this.err.String(), this.t)
	}
	return this
}

