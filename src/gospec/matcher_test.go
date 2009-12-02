// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"testing";
	"fmt";
)


type spyErrorLogger struct {
	failures int;
	lastMessage string;
}

func (log *spyErrorLogger) logError(message string) {
	log.failures++;
	log.lastMessage = message;
}


func Test__String_should_equal_string(t *testing.T) {
	log := new(spyErrorLogger);
	newMatcher("hotdog", log).Should.Equal("hotdog");
	assertEquals(0, log.failures, t);
	
	newMatcher("hotdog", log).Should.Equal("carrot");
	assertEquals(1, log.failures, t);
	assertEquals("Expected 'carrot' but was 'hotdog'", log.lastMessage, t);
}

func Test__String_should_not_equal_string(t *testing.T) {
	log := new(spyErrorLogger);
	newMatcher("hotdog", log).ShouldNot.Equal("carrot");
	assertEquals(0, log.failures, t);
	
	newMatcher("hotdog", log).ShouldNot.Equal("hotdog");
	assertEquals(1, log.failures, t);
	assertEquals("Did not expect 'hotdog' but was 'hotdog'", log.lastMessage, t);
}

func Test__Int_should_equal_int(t *testing.T) {
	log := new(spyErrorLogger);
	newMatcher(42, log).Should.Equal(42);
	assertEquals(0, log.failures, t);
	
	newMatcher(42, log).Should.Equal(13);
	assertEquals(1, log.failures, t);
	assertEquals("Expected '13' but was '42'", log.lastMessage, t);
}

func Test__Struct_should_equal_struct(t *testing.T) {
	log := new(spyErrorLogger);
	newMatcher(DummyStruct{42, 1}, log).Should.Equal(DummyStruct{42, 2});
	assertEquals(0, log.failures, t);
	
	newMatcher(DummyStruct{42, 1}, log).Should.Equal(DummyStruct{13, 2});
	assertEquals(1, log.failures, t);
	assertEquals("Expected 'DummyStruct13' but was 'DummyStruct42'", log.lastMessage, t);
}

func Test__Struct_pointer_should_equal_struct_pointer(t *testing.T) {
	log := new(spyErrorLogger);
	newMatcher(&DummyStruct{42, 1}, log).Should.Equal(&DummyStruct{42, 2});
	assertEquals(0, log.failures, t);
	
	newMatcher(&DummyStruct{42, 1}, log).Should.Equal(&DummyStruct{13, 2});
	assertEquals(1, log.failures, t);
	assertEquals("Expected 'DummyStruct13' but was 'DummyStruct42'", log.lastMessage, t);
}


type DummyStruct struct {
	value int;
	ignoredValue int;
}

func (this DummyStruct) Equals(other interface{}) bool {
	// TODO: figure out a better way to do equality, with less code to implement in client code
	if that, ok := other.(DummyStruct); ok {
		return this.equals(&that)
	}
	if that, ok := other.(*DummyStruct); ok {
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

