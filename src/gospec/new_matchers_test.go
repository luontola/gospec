// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/list"
	"fmt"
	"math"
	"nanospec"
	"os"
)


func MatcherMessagesSpec(c nanospec.Context) {
	spy := new(SpyErrorLogger)
	m := newMatcherAdapter(nil, spy)

	c.Specify("Positive expectation failures are reported with the positive message", func() {
		m.Expect(1, DummyEquals, 1)
		c.Expect(spy.LastError()).Equals("")

		m.Expect(1, DummyEquals, 2)
		c.Expect(spy.LastError()).Equals("1 should equal 2")
	})
	c.Specify("Negative expectation failures are reported with the negative message", func() {
		m.Expect(1, Not(DummyEquals), 2)
		c.Expect(spy.LastError()).Equals("")

		m.Expect(1, Not(DummyEquals), 1)
		c.Expect(spy.LastError()).Equals("1 should NOT equal 1")
	})
	c.Specify("Errors in expectations are reported with the error message", func() {
		m.Expect(666, DummyEquals, 1)
		c.Expect(spy.LastError()).Equals("Error: 666")
	})
}

func DummyEquals(actual interface{}, expected interface{}) (match bool, pos os.Error, neg os.Error, err os.Error) {
	if actual.(int) == 666 {
		err = Errorf("Error: %v", actual)
		return
	}
	match = actual == expected
	pos = Errorf("%v should equal %v", actual, expected)
	neg = Errorf("%v should NOT equal %v", actual, expected)
	return
}


func MatchersSpec(c nanospec.Context) {

	c.Specify("Matcher: Equals", func() {
		c.Specify("strings", func() {
			expectation(c, "apple", Equals, "apple").Passes()
			expectation(c, "apple", Equals, "orange").Fails().
				WithMessage(
				"Expected 'orange' but was 'apple'",
				"Did not expect 'orange' but was 'apple'")
		})
		c.Specify("ints", func() {
			expectation(c, 42, Equals, 42).Passes()
			expectation(c, 42, Equals, 999).Fails()
		})
		c.Specify("structs", func() {
			expectation(c, DummyStruct{42, 1}, Equals, DummyStruct{42, 2}).Passes()
			expectation(c, DummyStruct{42, 1}, Equals, DummyStruct{999, 2}).Fails()
		})
		c.Specify("struct pointers", func() {
			expectation(c, &DummyStruct{42, 1}, Equals, &DummyStruct{42, 2}).Passes()
			expectation(c, &DummyStruct{42, 1}, Equals, &DummyStruct{999, 2}).Fails()
		})
	})

	c.Specify("Matcher: IsSame", func() {
		a1 := new(os.File)
		a2 := a1
		b := new(os.File)
		expectation(c, a1, IsSame, a2).Passes()
		expectation(c, a1, IsSame, b).Fails().
			WithMessage(
			fmt.Sprintf("Expected '%v' but was '%v'", b, a1),
			fmt.Sprintf("Did not expect '%v' but was '%v'", b, a1))
		expectation(c, 1, IsSame, b).GivesError("Expected a pointer, but was '1' of type 'int'")
		expectation(c, b, IsSame, 1).GivesError("Expected a pointer, but was '1' of type 'int'")
	})

	c.Specify("Matcher: IsNil", func() {
		expectation(c, nil, IsNil).Passes()         // interface value nil
		expectation(c, (*int)(nil), IsNil).Passes() // typed pointer nil inside an interface value
		expectation(c, new(int), IsNil).Fails()
		expectation(c, 1, IsNil).Fails().
			WithMessage(
			"Expected <nil> but was '1'",
			"Did not expect <nil> but was '1'")
	})

	c.Specify("Matcher: IsTrue", func() {
		expectation(c, true, IsTrue).Passes()
		expectation(c, false, IsTrue).Fails().
			WithMessage(
			"Expected 'true' but was 'false'",
			"Did not expect 'true' but was 'false'")
	})

	c.Specify("Matcher: IsFalse", func() {
		expectation(c, false, IsFalse).Passes()
		expectation(c, true, IsFalse).Fails().
			WithMessage(
			"Expected 'false' but was 'true'",
			"Did not expect 'false' but was 'true'")
	})

	c.Specify("Matcher: Satisfies", func() {
		value := 42
		expectation(c, value, Satisfies, value < 100).Passes()
		expectation(c, value, Satisfies, value > 100).Fails().
			WithMessage(
			"Criteria not satisfied by '42'",
			"Criteria not satisfied by '42'")
	})

	c.Specify("Matcher: IsWithin", func() {
		value := float64(3.141)
		pi := float64(math.Pi)
		expectation(c, value, IsWithin(0.001), pi).Passes()
		expectation(c, value, IsWithin(0.0001), pi).Fails().
			WithMessage(
			"Expected '3.141592653589793' ± 0.0001 but was '3.141'",
			"Did not expect '3.141592653589793' ± 0.0001 but was '3.141'")

		c.Specify("cannot compare ints", func() {
			value := int(3)
			pi := float64(math.Pi)
			expectation(c, value, IsWithin(0.001), pi).GivesError("Expected a float, but was '3' of type 'int'")
			expectation(c, pi, IsWithin(0.001), value).GivesError("Expected a float, but was '3' of type 'int'")
		})
	})

	c.Specify("Matcher: Contains", func() {
		values := []string{"one", "two", "three"}

		expectation(c, values, Contains, "one").Passes()
		expectation(c, values, Contains, "two").Passes()
		expectation(c, values, Contains, "three").Passes()

		expectation(c, values, Contains, "four").Fails().
			WithMessage(
			"Expected 'four' to be in '[one two three]' but it was not",
			"Did not expect 'four' to be in '[one two three]' but it was")
	})

	c.Specify("Matcher: ContainsAll", func() {
		values := []string{"one", "two", "three"}

		expectation(c, values, ContainsAll, Values()).Passes()
		expectation(c, values, ContainsAll, Values("one")).Passes()
		expectation(c, values, ContainsAll, Values("three", "two")).Passes()
		expectation(c, values, ContainsAll, Values("one", "two", "three")).Passes()

		expectation(c, values, ContainsAll, Values("four")).Fails()
		expectation(c, values, ContainsAll, Values("one", "four")).Fails().
			WithMessage(
			"Expected all of '[one four]' to be in '[one two three]' but they were not",
			"Did not expect all of '[one four]' to be in '[one two three]' but they were")
	})

	c.Specify("Matcher: ContainsAny", func() {
		values := []string{"one", "two", "three"}

		expectation(c, values, ContainsAny, Values("one")).Passes()
		expectation(c, values, ContainsAny, Values("three", "two")).Passes()
		expectation(c, values, ContainsAny, Values("four", "one", "five")).Passes()
		expectation(c, values, ContainsAny, Values("one", "two", "three")).Passes()

		expectation(c, values, ContainsAny, Values()).Fails()
		expectation(c, values, ContainsAny, Values("four")).Fails()
		expectation(c, values, ContainsAny, Values("four", "five")).Fails().
			WithMessage(
			"Expected any of '[four five]' to be in '[one two three]' but they were not",
			"Did not expect any of '[four five]' to be in '[one two three]' but they were")
	})

	c.Specify("Matcher: ContainsExactly", func() {
		values := []string{"one", "two", "three"}

		expectation(c, values, ContainsExactly, Values("one", "two", "three")).Passes()
		expectation(c, values, ContainsExactly, Values("three", "one", "two")).Passes()

		expectation(c, values, ContainsExactly, Values()).Fails()
		expectation(c, values, ContainsExactly, Values("four")).Fails()
		expectation(c, values, ContainsExactly, Values("one", "two")).Fails()
		expectation(c, values, ContainsExactly, Values("one", "two", "three", "four")).Fails().
			WithMessage(
			"Expected exactly '[one two three four]' to be in '[one two three]' but they were not",
			"Did not expect exactly '[one two three four]' to be in '[one two three]' but they were")

		// duplicate values are allowed
		values = []string{"a", "a", "b"}

		expectation(c, values, ContainsExactly, Values("a", "a", "b")).Passes()
		expectation(c, values, ContainsExactly, Values("a", "b", "a")).Passes()

		expectation(c, values, ContainsExactly, Values("a", "b", "b")).Fails()
		expectation(c, values, ContainsExactly, Values("a", "a", "a", "b")).Fails()
		expectation(c, values, ContainsExactly, Values("a", "a", "b", "b")).Fails()
	})

	c.Specify("Matcher: ContainsInOrder", func() {
		values := []string{"one", "two", "three"}

		expectation(c, values, ContainsInOrder, Values("one", "two", "three")).Passes()

		expectation(c, values, ContainsInOrder, Values()).Fails()
		expectation(c, values, ContainsInOrder, Values("one", "two")).Fails()
		expectation(c, values, ContainsInOrder, Values("one", "two", "four")).Fails()
		expectation(c, values, ContainsInOrder, Values("one", "two", "three", "four")).Fails()
		expectation(c, values, ContainsInOrder, Values("three", "one", "two")).Fails().
			WithMessage(
			"Expected in order '[three one two]' to be in '[one two three]' but they were not",
			"Did not expect in order '[three one two]' to be in '[one two three]' but they were")
	})

	c.Specify("Matcher: ContainsInPartialOrder", func() {
		values := []string{"1", "2", "2", "3", "4"}

		expectation(c, values, ContainsInPartialOrder, Values()).Passes()
		expectation(c, values, ContainsInPartialOrder, Values("1")).Passes()
		expectation(c, values, ContainsInPartialOrder, Values("1", "2", "2")).Passes()
		expectation(c, values, ContainsInPartialOrder, Values("1", "2", "3")).Passes()
		expectation(c, values, ContainsInPartialOrder, Values("1", "2", "2", "3", "4")).Passes()

		expectation(c, values, ContainsInPartialOrder, Values("1", "1")).Fails()
		expectation(c, values, ContainsInPartialOrder, Values("2", "1")).Fails()
		expectation(c, values, ContainsInPartialOrder, Values("2", "2", "2")).Fails()
		expectation(c, values, ContainsInPartialOrder, Values("1", "4", "3")).Fails().
			WithMessage(
			"Expected in partial order '[1 4 3]' to be in '[1 2 2 3 4]' but they were not",
			"Did not expect in partial order '[1 4 3]' to be in '[1 2 2 3 4]' but they were")
	})

	c.Specify("Conversions for containment matchers", func() {
		c.Specify("array to array", func() {
			values := [...]string{"one", "two", "three"}

			result, _ := toArray(values)

			c.Expect(len(result)).Equals(3)
			c.Expect(result[0]).Equals("one")
			c.Expect(result[1]).Equals("two")
			c.Expect(result[2]).Equals("three")
		})
		c.Specify("channel to array", func() {
			values := list.New()
			values.PushBack("one")
			values.PushBack("two")
			values.PushBack("three")

			result, _ := toArray(values.Iter())

			c.Expect(len(result)).Equals(3)
			c.Expect(result[0]).Equals("one")
			c.Expect(result[1]).Equals("two")
			c.Expect(result[2]).Equals("three")
		})
		c.Specify("iterable to array", func() {
			values := list.New()
			values.PushBack("one")
			values.PushBack("two")
			values.PushBack("three")

			result, _ := toArray(values)

			c.Expect(len(result)).Equals(3)
			c.Expect(result[0]).Equals("one")
			c.Expect(result[1]).Equals("two")
			c.Expect(result[2]).Equals("three")
		})
		c.Specify("unsupported value to array", func() {
			_, err := toArray("foo")
			c.Expect(err.String()).Equals("Unknown type 'string', not iterable: foo")
		})
	})
}


// Used by the Equals matcher's tests
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


// Test utilities

func expectation(c nanospec.Context, actual interface{}, matcher Matcher, expected ...interface{}) *matchAssert {
	match, pos, neg, err := matcher.Match(actual, expected)
	return &matchAssert{match, pos, neg, err, c}
}

type matchAssert struct {
	match bool
	pos   os.Error
	neg   os.Error
	err   os.Error
	c     nanospec.Context
}

func (this *matchAssert) Passes() *matchAssert {
	if this.err != nil {
		this.c.Errorf("expected to pass, but had an error: %v", this.err)
	}
	if !this.match {
		this.c.Errorf("expected to pass, but failed")
	}
	return this
}

func (this *matchAssert) Fails() *matchAssert {
	if this.err != nil {
		this.c.Errorf("expected to fail, but had an error: %v", this.err)
	}
	if this.match {
		this.c.Errorf("expected to fail, but passed")
	}
	return this
}

func (this *matchAssert) WithMessage(expectedPos string, expectedNeg string) *matchAssert {
	this.c.Expect(this.pos.String()).Equals(expectedPos)
	this.c.Expect(this.neg.String()).Equals(expectedNeg)
	return this
}

func (this *matchAssert) GivesError(expectedErr string) *matchAssert {
	if this.err == nil {
		this.c.Errorf("expected have an error, but did not")
	} else {
		this.c.Expect(this.err.String()).Equals(expectedErr)
	}
	return this
}


type SpyErrorLogger struct {
	lastError *Error
}

func (this *SpyErrorLogger) AddError(error *Error) {
	this.lastError = error
}

func (this *SpyErrorLogger) AddFatalError(error *Error) {
	this.AddError(error)
}

func (this *SpyErrorLogger) LastError() string {
	defer this.Reset()
	if this.lastError == nil {
		return ""
	}
	return this.lastError.Message
}

func (this *SpyErrorLogger) Reset() {
	this.lastError = nil
}
