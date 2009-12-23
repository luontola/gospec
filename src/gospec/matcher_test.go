// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/list"
	"testing"
	"fmt"
)


// "Equal" matcher

func Test__String_should_EQUAL_string(t *testing.T) {
	log := new(spyErrorLogger)

	newMatcher("apple", log).Should.Equal("apple")
	log.ShouldHaveNoErrors(t)

	newMatcher("apple", log).Should.Equal("orange")
	log.ShouldHaveTheError("Expected 'orange' but was 'apple'", t)
}

func Test__String_should_NOT_EQUAL_string(t *testing.T) {
	log := new(spyErrorLogger)

	newMatcher("apple", log).ShouldNot.Equal("orange")
	log.ShouldHaveNoErrors(t)

	newMatcher("apple", log).ShouldNot.Equal("apple")
	log.ShouldHaveTheError("Did not expect 'apple' but was 'apple'", t)
}

func Test__Int_should_EQUAL_int(t *testing.T) {
	log := new(spyErrorLogger)

	newMatcher(42, log).Should.Equal(42)
	log.ShouldHaveNoErrors(t)

	newMatcher(42, log).Should.Equal(999)
	log.ShouldHaveTheError("Expected '999' but was '42'", t)
}

func Test__Struct_should_EQUAL_struct(t *testing.T) {
	log := new(spyErrorLogger)

	newMatcher(DummyStruct{42, 1}, log).Should.Equal(DummyStruct{42, 2})
	log.ShouldHaveNoErrors(t)

	newMatcher(DummyStruct{42, 1}, log).Should.Equal(DummyStruct{999, 2})
	log.ShouldHaveTheError("Expected 'DummyStruct999' but was 'DummyStruct42'", t)
}

func Test__Struct_pointer_should_EQUAL_struct_pointer(t *testing.T) {
	log := new(spyErrorLogger)

	newMatcher(&DummyStruct{42, 1}, log).Should.Equal(&DummyStruct{42, 2})
	log.ShouldHaveNoErrors(t)

	newMatcher(&DummyStruct{42, 1}, log).Should.Equal(&DummyStruct{999, 2})
	log.ShouldHaveTheError("Expected 'DummyStruct999' but was 'DummyStruct42'", t)
}


// "Be" matcher

func Test__Object_should_BE_some_expression(t *testing.T) {
	log := new(spyErrorLogger)
	value := 42
	
	newMatcher(value, log).Should.Be(value > 40)
	log.ShouldHaveNoErrors(t)
	
	newMatcher(value, log).Should.Be(value > 999)
	log.ShouldHaveTheError("Criteria not satisfied by '42'", t)
}

func Test__Object_should_NOT_BE_some_expression(t *testing.T) {
	log := new(spyErrorLogger)
	value := 42
	
	newMatcher(value, log).ShouldNot.Be(value < 40)
	log.ShouldHaveNoErrors(t)
	
	newMatcher(value, log).ShouldNot.Be(value < 999)
	log.ShouldHaveTheError("Criteria not satisfied by '42'", t)
}


// "Contain" matcher

func Test__Array_should_CONTAIN_a_value(t *testing.T) {
	log := new(spyErrorLogger)
	values := [...]string{"one", "two", "three"}
	
	newMatcher(values, log).Should.Contain("one")
	newMatcher(values, log).Should.Contain("two")
	newMatcher(values, log).Should.Contain("three")
	log.ShouldHaveNoErrors(t)
	
	newMatcher(values, log).Should.Contain("four")
	log.ShouldHaveTheError("Expected 'four' to be in '[one two three]' but it was not", t)
}

func Test__Array_should_NOT_CONTAIN_a_value(t *testing.T) {
	log := new(spyErrorLogger)
	values := [...]string{"one", "two", "three"}
	
	newMatcher(values, log).ShouldNot.Contain("four")
	log.ShouldHaveNoErrors(t)
	
	newMatcher(values, log).ShouldNot.Contain("one")
	log.ShouldHaveTheError("Did not expect 'one' to be in '[one two three]' but it was", t)
}

func Test__Iterable_should_CONTAIN_a_value(t *testing.T) {
	log := new(spyErrorLogger)
	values := list.New()
	values.PushBack("one")
	values.PushBack("two")
	values.PushBack("three")
	
	newMatcher(values.Iter(), log).Should.Contain("one")
	newMatcher(values.Iter(), log).Should.Contain("two")
	newMatcher(values.Iter(), log).Should.Contain("three")
	log.ShouldHaveNoErrors(t)
	
	newMatcher(values.Iter(), log).Should.Contain("four")
	log.ShouldHaveTheError("Expected 'four' to be in '[one two three]' but it was not", t)
}

func Test__Iterable_should_NOT_CONTAIN_a_value(t *testing.T) {
	log := new(spyErrorLogger)
	values := list.New()
	values.PushBack("one")
	values.PushBack("two")
	values.PushBack("three")
	
	newMatcher(values.Iter(), log).ShouldNot.Contain("four")
	log.ShouldHaveNoErrors(t)
	
	newMatcher(values.Iter(), log).ShouldNot.Contain("one")
	log.ShouldHaveTheError("Did not expect 'one' to be in '[one two three]' but it was", t)
}


// Utilities

type spyErrorLogger struct {
	failures    int
	lastMessage string
}

func (log *spyErrorLogger) AddError(message string) {
	log.failures++
	log.lastMessage = message
}

func (log *spyErrorLogger) AddFatalError(message string) {
	log.AddError(message)
}

func (log *spyErrorLogger) Reset() {
	log.failures = 0
	log.lastMessage = ""
}

func (log *spyErrorLogger) ShouldHaveNoErrors(t *testing.T) {
	assertEquals(0, log.failures, t)
	log.Reset()
}

func (log *spyErrorLogger) ShouldHaveTheError(message string, t *testing.T) {
	assertEquals(1, log.failures, t)
	assertEquals(message, log.lastMessage, t)
	log.Reset()
}


// Dummies

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

