// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"testing"
	"fmt"
)


func Test__String_should_equal_string(t *testing.T) {
	log := new(spyErrorLogger)

	newMatcher("hotdog", log).Should.Equal("hotdog")
	log.ShouldHaveNoErrors(t)

	newMatcher("hotdog", log).Should.Equal("carrot")
	log.ShouldHaveTheError("Expected 'carrot' but was 'hotdog'", t)
}

func Test__String_should_not_equal_string(t *testing.T) {
	log := new(spyErrorLogger)

	newMatcher("hotdog", log).ShouldNot.Equal("carrot")
	log.ShouldHaveNoErrors(t)

	newMatcher("hotdog", log).ShouldNot.Equal("hotdog")
	log.ShouldHaveTheError("Did not expect 'hotdog' but was 'hotdog'", t)
}

func Test__Int_should_equal_int(t *testing.T) {
	log := new(spyErrorLogger)

	newMatcher(42, log).Should.Equal(42)
	log.ShouldHaveNoErrors(t)

	newMatcher(42, log).Should.Equal(13)
	log.ShouldHaveTheError("Expected '13' but was '42'", t)
}

func Test__Struct_should_equal_struct(t *testing.T) {
	log := new(spyErrorLogger)

	newMatcher(DummyStruct{42, 1}, log).Should.Equal(DummyStruct{42, 2})
	log.ShouldHaveNoErrors(t)

	newMatcher(DummyStruct{42, 1}, log).Should.Equal(DummyStruct{13, 2})
	log.ShouldHaveTheError("Expected 'DummyStruct13' but was 'DummyStruct42'", t)
}

func Test__Struct_pointer_should_equal_struct_pointer(t *testing.T) {
	log := new(spyErrorLogger)

	newMatcher(&DummyStruct{42, 1}, log).Should.Equal(&DummyStruct{42, 2})
	log.ShouldHaveNoErrors(t)

	newMatcher(&DummyStruct{42, 1}, log).Should.Equal(&DummyStruct{13, 2})
	log.ShouldHaveTheError("Expected 'DummyStruct13' but was 'DummyStruct42'", t)
}


type spyErrorLogger struct {
	failures    int
	lastMessage string
}

func (log *spyErrorLogger) AddError(message string) {
	log.failures++
	log.lastMessage = message
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

