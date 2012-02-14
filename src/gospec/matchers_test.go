// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/list"
	"errors"
	"fmt"
	"math"
	"github.com/orfjackal/nanospec.go/src/nanospec"
	"os"
)

func MatcherMessagesSpec(c nanospec.Context) {
	spy := new(SpyErrorLogger)
	m := newMatcherAdapter(nil, spy, ExpectFailed)

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
		c.Expect(spy.LastError()).Equals("666 illegal value")
	})
}

func DummyEquals(actual interface{}, expected interface{}) (match bool, pos Message, neg Message, err error) {
	if actual.(int) == 666 {
		err = Errorf("illegal value")
		return
	}
	match = actual == expected
	pos = Messagef(actual, "should equal %v", expected)
	neg = Messagef(actual, "should NOT equal %v", expected)
	return
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
	return this.lastError.Actual + " " + this.lastError.Message
}

func (this *SpyErrorLogger) Reset() {
	this.lastError = nil
}

func MatchersSpec(c nanospec.Context) {

	c.Specify("Matcher: Equals", func() {
		c.Specify("strings", func() {
			c.Expect(E("apple", Equals, "apple")).Matches(Passes)
			c.Expect(E("apple", Equals, "orange")).Matches(FailsWithMessage(
				"equals “orange”",
				"does NOT equal “orange”"))
		})
		c.Specify("ints", func() {
			c.Expect(E(42, Equals, 42)).Matches(Passes)
			c.Expect(E(42, Equals, 999)).Matches(Fails)
		})
		c.Specify("structs", func() {
			c.Expect(E(DummyStruct{42, 1}, Equals, DummyStruct{42, 2})).Matches(Passes)
			c.Expect(E(DummyStruct{42, 1}, Equals, DummyStruct{999, 2})).Matches(Fails)
		})
		c.Specify("struct pointers", func() {
			c.Expect(E(&DummyStruct{42, 1}, Equals, &DummyStruct{42, 2})).Matches(Passes)
			c.Expect(E(&DummyStruct{42, 1}, Equals, &DummyStruct{999, 2})).Matches(Fails)
		})
	})

	c.Specify("Matcher: IsSame", func() {
		a1 := new(os.File)
		a2 := a1
		b := new(os.File)

		c.Expect(E(a1, IsSame, a2)).Matches(Passes)
		c.Expect(E(a1, IsSame, b)).Matches(FailsWithMessage(
			fmt.Sprintf("is same as “%v”", b),
			fmt.Sprintf("is NOT same as “%v”", b)))

		c.Specify("cannot compare values, but only pointers", func() {
			c.Expect(E(1, IsSame, b)).Matches(GivesError("type error: expected a pointer, but was “1” of type “int”"))
			c.Expect(E(b, IsSame, 1)).Matches(GivesError("type error: expected a pointer, but was “1” of type “int”"))
		})
	})

	c.Specify("Matcher: IsNil", func() {
		c.Expect(E(nil, IsNil)).Matches(Passes)         // interface value nil
		c.Expect(E((*int)(nil), IsNil)).Matches(Passes) // typed pointer nil inside an interface value
		c.Expect(E(new(int), IsNil)).Matches(Fails)
		c.Expect(E(1, IsNil)).Matches(FailsWithMessage(
			"is <nil>",
			"is NOT <nil>"))
	})

	c.Specify("Matcher: IsTrue", func() {
		c.Expect(E(true, IsTrue)).Matches(Passes)
		c.Expect(E(false, IsTrue)).Matches(FailsWithMessage(
			"is <true>",
			"is NOT <true>"))
	})

	c.Specify("Matcher: IsFalse", func() {
		c.Expect(E(false, IsFalse)).Matches(Passes)
		c.Expect(E(true, IsFalse)).Matches(FailsWithMessage(
			"is <false>",
			"is NOT <false>"))
	})

	c.Specify("Matcher: Satisfies", func() {
		value := 42

		c.Expect(E(value, Satisfies, value < 100)).Matches(Passes)
		c.Expect(E(value, Satisfies, value > 100)).Matches(FailsWithMessage(
			"satisfies the criteria",
			"does NOT satisfy the criteria"))
	})

	c.Specify("Matcher: IsWithin", func() {
		value := float64(3.141)
		pi := float64(math.Pi)

		c.Expect(E(value, IsWithin(0.001), pi)).Matches(Passes)
		c.Expect(E(value, IsWithin(0.0001), pi)).Matches(FailsWithMessage(
			"is within 3.141592653589793 ± 0.0001",
			"is NOT within 3.141592653589793 ± 0.0001"))

		c.Specify("cannot compare ints", func() {
			value := int(3)
			pi := float64(math.Pi)

			c.Expect(E(value, IsWithin(0.001), pi)).Matches(GivesError("type error: expected a float, but was “3” of type “int”"))
			c.Expect(E(pi, IsWithin(0.001), value)).Matches(GivesError("type error: expected a float, but was “3” of type “int”"))
		})
	})

	c.Specify("Matcher: Contains", func() {
		values := []string{"one", "two", "three"}

		c.Expect(E(values, Contains, "one")).Matches(Passes)
		c.Expect(E(values, Contains, "two")).Matches(Passes)
		c.Expect(E(values, Contains, "three")).Matches(Passes)

		c.Expect(E(values, Contains, "four")).Matches(FailsWithMessage(
			"contains “four”",
			"does NOT contain “four”"))
	})

	c.Specify("Matcher: ContainsAll", func() {
		values := []string{"one", "two", "three"}

		c.Expect(E(values, ContainsAll, Values())).Matches(Passes)
		c.Expect(E(values, ContainsAll, Values("one"))).Matches(Passes)
		c.Expect(E(values, ContainsAll, Values("three", "two"))).Matches(Passes)
		c.Expect(E(values, ContainsAll, Values("one", "two", "three"))).Matches(Passes)

		c.Expect(E(values, ContainsAll, Values("four"))).Matches(Fails)
		c.Expect(E(values, ContainsAll, Values("one", "four"))).Matches(FailsWithMessage(
			"contains all of “[one four]”",
			"does NOT contain all of “[one four]”"))
	})

	c.Specify("Matcher: ContainsAny", func() {
		values := []string{"one", "two", "three"}

		c.Expect(E(values, ContainsAny, Values("one"))).Matches(Passes)
		c.Expect(E(values, ContainsAny, Values("three", "two"))).Matches(Passes)
		c.Expect(E(values, ContainsAny, Values("four", "one", "five"))).Matches(Passes)
		c.Expect(E(values, ContainsAny, Values("one", "two", "three"))).Matches(Passes)

		c.Expect(E(values, ContainsAny, Values())).Matches(Fails)
		c.Expect(E(values, ContainsAny, Values("four"))).Matches(Fails)
		c.Expect(E(values, ContainsAny, Values("four", "five"))).Matches(FailsWithMessage(
			"contains any of “[four five]”",
			"does NOT contain any of “[four five]”"))
	})

	c.Specify("Matcher: ContainsExactly", func() {
		values := []string{"one", "two", "three"}

		c.Expect(E(values, ContainsExactly, Values("one", "two", "three"))).Matches(Passes)
		c.Expect(E(values, ContainsExactly, Values("three", "one", "two"))).Matches(Passes)

		c.Expect(E(values, ContainsExactly, Values())).Matches(Fails)
		c.Expect(E(values, ContainsExactly, Values("four"))).Matches(Fails)
		c.Expect(E(values, ContainsExactly, Values("one", "two"))).Matches(Fails)
		c.Expect(E(values, ContainsExactly, Values("one", "two", "three", "four"))).Matches(FailsWithMessage(
			"contains exactly “[one two three four]”",
			"does NOT contain exactly “[one two three four]”"))

		// duplicate values are allowed
		values = []string{"a", "a", "b"}

		c.Expect(E(values, ContainsExactly, Values("a", "a", "b"))).Matches(Passes)
		c.Expect(E(values, ContainsExactly, Values("a", "b", "a"))).Matches(Passes)

		c.Expect(E(values, ContainsExactly, Values("a", "b", "b"))).Matches(Fails)
		c.Expect(E(values, ContainsExactly, Values("a", "a", "a", "b"))).Matches(Fails)
		c.Expect(E(values, ContainsExactly, Values("a", "a", "b", "b"))).Matches(Fails)
	})

	c.Specify("Matcher: ContainsInOrder", func() {
		values := []string{"one", "two", "three"}

		c.Expect(E(values, ContainsInOrder, Values("one", "two", "three"))).Matches(Passes)

		c.Expect(E(values, ContainsInOrder, Values())).Matches(Fails)
		c.Expect(E(values, ContainsInOrder, Values("one", "two"))).Matches(Fails)
		c.Expect(E(values, ContainsInOrder, Values("one", "two", "four"))).Matches(Fails)
		c.Expect(E(values, ContainsInOrder, Values("one", "two", "three", "four"))).Matches(Fails)
		c.Expect(E(values, ContainsInOrder, Values("three", "one", "two"))).Matches(FailsWithMessage(
			"contains in order “[three one two]”",
			"does NOT contain in order “[three one two]”"))
	})

	c.Specify("Matcher: ContainsInPartialOrder", func() {
		values := []string{"1", "2", "2", "3", "4"}

		c.Expect(E(values, ContainsInPartialOrder, Values())).Matches(Passes)
		c.Expect(E(values, ContainsInPartialOrder, Values("1"))).Matches(Passes)
		c.Expect(E(values, ContainsInPartialOrder, Values("1", "2", "2"))).Matches(Passes)
		c.Expect(E(values, ContainsInPartialOrder, Values("1", "2", "3"))).Matches(Passes)
		c.Expect(E(values, ContainsInPartialOrder, Values("1", "2", "2", "3", "4"))).Matches(Passes)

		c.Expect(E(values, ContainsInPartialOrder, Values("1", "1"))).Matches(Fails)
		c.Expect(E(values, ContainsInPartialOrder, Values("2", "1"))).Matches(Fails)
		c.Expect(E(values, ContainsInPartialOrder, Values("2", "2", "2"))).Matches(Fails)
		c.Expect(E(values, ContainsInPartialOrder, Values("1", "4", "3"))).Matches(FailsWithMessage(
			"contains in partial order “[1 4 3]”",
			"does NOT contain in partial order “[1 4 3]”"))
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
			values := make(chan string, 10)
			values <- "one"
			values <- "two"
			values <- "three"
			close(values)

			result, _ := toArray(values)

			c.Expect(len(result)).Equals(3)
			c.Expect(result[0]).Equals("one")
			c.Expect(result[1]).Equals("two")
			c.Expect(result[2]).Equals("three")
		})
		c.Specify("list to array", func() {
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
			c.Expect(err.Error()).Equals("type error: expected a collection type, but was “foo” of type “string”")
		})
	})

	c.Specify("Containment matchers convert the actual value to a printable array "+
		"when it's not originally printable",
		func() {
			values := list.New()
			values.PushBack("one")
			values.PushBack("two")
			values.PushBack("three")
			nonPrintableList := values
			listAsString := "[one two three]"
			c.Expect(fmt.Sprint(nonPrintableList)).NotEquals(listAsString)

			c.Specify("Matcher: Contains", func() {
				c.Expect(func(c Context) {
					c.Expect(nonPrintableList, Contains, "xxx")
				}).Matches(SpecsReportContains(listAsString))
			})

			multiValueMatchers := []Matcher{
				ContainsAll,
				ContainsAny,
				ContainsExactly,
				ContainsInOrder,
				ContainsInPartialOrder,
			}
			for _, matcher := range multiValueMatchers {
				c.Specify("Matcher: "+functionName(matcher), func() {
					c.Expect(func(c Context) {
						c.Expect(nonPrintableList, matcher, Values("xxx"))
					}).Matches(SpecsReportContains(listAsString))
				})
			}
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

type ExpectationHolder struct {
	match bool
	pos   Message
	neg   Message
	err   error
}

func E(actual interface{}, matcher Matcher, expected ...interface{}) *ExpectationHolder {
	match, pos, neg, err := matcher.Match(actual, expected...)
	return &ExpectationHolder{match, pos, neg, err}
}

func Passes(v interface{}) error {
	ex := v.(*ExpectationHolder)
	if ex.match && ex.err == nil {
		return nil
	}
	return ex.ToError()
}

func Fails(v interface{}) error {
	ex := v.(*ExpectationHolder)
	if !ex.match && ex.err == nil {
		return nil
	}
	return ex.ToError()
}

func FailsWithMessage(pos string, neg string) nanospec.Matcher {
	return func(v interface{}) error {
		ex := v.(*ExpectationHolder)
		if !ex.match &&
			ex.pos != nil &&
			ex.pos.Expectation() == pos &&
			ex.neg != nil &&
			ex.neg.Expectation() == neg &&
			ex.err == nil {
			return nil
		}
		return ex.ToError()
	}
}

func GivesError(err string) nanospec.Matcher {
	return func(v interface{}) error {
		ex := v.(*ExpectationHolder)
		if ex.err != nil && ex.err.Error() == err {
			return nil
		}
		return ex.ToError()
	}
}

func (this *ExpectationHolder) ToError() error {
	return errors.New(fmt.Sprintf(
		"Mather failed its expectations\n\tmatch: %v\n\tpos: %v\n\tneg: %v\n\terr: %v",
		this.match, this.pos, this.neg, this.err))
}
