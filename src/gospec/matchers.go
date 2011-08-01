// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt"
	"container/vector"
	"container/list"
	"math"
	"os"
	"reflect"
)


type matcherAdapter struct {
	location    *Location
	log         errorLogger
	matcherType ErrorType
}

func newMatcherAdapter(location *Location, log errorLogger, matcherType ErrorType) *matcherAdapter {
	return &matcherAdapter{location, log, matcherType}
}

func (this *matcherAdapter) Expect(actual interface{}, matcher Matcher, expected ...interface{}) {
	match, pos, _, err := matcher.Match(actual, expected...)
	if err != nil {
		this.addError(err, actual)
	} else if !match {
		this.addFailure(pos)
	}
}

func (this *matcherAdapter) addFailure(message Message) {
	this.writeToLog(this.matcherType, message.Expectation(), message.Actual())
}

func (this *matcherAdapter) addError(err os.Error, actual interface{}) {
	this.writeToLog(OtherError, err.String(), actual)
}

func (this *matcherAdapter) writeToLog(errortype ErrorType, message string, actual interface{}) {
	stacktrace := toStackTrace(this.location)
	e := newError(errortype, message, fmt.Sprint(actual), stacktrace)
	this.log.AddError(e)
}

func toStackTrace(loc *Location) []*Location {
	if loc != nil {
		return []*Location{loc}
	}
	return []*Location{}
}


// Matchers are used in expectations to compare the actual and expected values.
//
// Return values:
//   match: Should be true when `actual` and `expected` match, otherwise false.
//   pos:   Message for a failed expectation.
//   neg:   Message for a failed expectation when the matcher is combined with Not.
//   err:   Message for an unrecoverable error, for example if the arguments had a wrong type.
type Matcher func(actual interface{}, expected interface{}) (match bool, pos Message, neg Message, err os.Error)

// Calls the matcher with the actual value and an optional expected value.
// If no expected value is given, then <nil> will be used.
func (matcher Matcher) Match(actual interface{}, optionalExpected ...interface{}) (match bool, pos Message, neg Message, err os.Error) {
	var expected interface{}
	if len(optionalExpected) > 0 {
		expected = optionalExpected[0]
	}
	match, pos, neg, err = matcher(actual, expected)
	return
}


type Message interface {
	Actual() interface{}
	Expectation() string
}

func Messagef(actual interface{}, expectationFormat string, expectationArgs ...interface{}) Message {
	expectation := Errorf(expectationFormat, expectationArgs...)
	return &message{actual, expectation}
}

type message struct {
	actual      interface{}
	expectation os.Error
}

func (this *message) Actual() interface{} {
	return this.actual
}

func (this *message) Expectation() string {
	return this.expectation.String()
}


// Constructs an error message the same way as fmt.Sprintf(), but the string is
// created lazily when it is used, if it is used at all. This avoids unnecessary
// string parsing in matchers, because most of the time there are no failures
// and thus the error messages are not used.
func Errorf(format string, args ...interface{}) os.Error {
	return lazyString(func() string {
		return fmt.Sprintf(format, args...)
	})
}

type lazyString func() string

func (this lazyString) String() string {
	return this()
}


// Easy array creation, to give multiple expected values to a matcher.
func Values(values ...interface{}) []interface{} {
	return values
}


// Negates the meaning of a Matcher. Matches when the original matcher does not
// match, and the other way around.
func Not(matcher Matcher) Matcher {
	return func(actual interface{}, expected interface{}) (match bool, pos Message, neg Message, err os.Error) {
		match, pos, neg, err = matcher(actual, expected)
		match = !match
		pos, neg = neg, pos
		return
	}
}


// The actual value must equal the expected value. For primitives the equality
// operator is used. All other objects must implement the Equality interface.
func Equals(actual interface{}, expected interface{}) (match bool, pos Message, neg Message, err os.Error) {
	match = areEqual(actual, expected)
	pos = Messagef(actual, "equals “%v”", expected)
	neg = Messagef(actual, "does NOT equal “%v”", expected)
	return
}

func areEqual(a interface{}, b interface{}) bool {
	if a2, ok := a.(Equality); ok {
		return a2.Equals(b)
	}
	return a == b
}

type Equality interface {
	Equals(other interface{}) bool
}


// The actual value must be a pointer to the same object as the expected value.
func IsSame(actual interface{}, expected interface{}) (match bool, pos Message, neg Message, err os.Error) {
	ptr1, err := pointerOf(actual)
	if err != nil {
		return
	}
	ptr2, err := pointerOf(expected)
	if err != nil {
		return
	}
	match = ptr1 == ptr2
	pos = Messagef(actual, "is same as “%v”", expected)
	neg = Messagef(actual, "is NOT same as “%v”", expected)
	return
}

func pointerOf(value interface{}) (ptr uintptr, err os.Error) {
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Ptr:
		ptr = v.Pointer()
	default:
		err = Errorf("type error: expected a pointer, but was “%v” of type “%T”", value, value)
	}
	return
}


// The actual value must be <nil>, or a typed nil pointer inside an interface value.
// See http://groups.google.com/group/golang-nuts/browse_thread/thread/d900674d491ef8d
// for discussion on how in Go typed nil values can turn into non-nil interface values.
func IsNil(actual interface{}, _ interface{}) (match bool, pos Message, neg Message, err os.Error) {
	match = actual == nil || isNilPointerInsideInterfaceValue(actual)
	pos = Messagef(actual, "is <nil>")
	neg = Messagef(actual, "is NOT <nil>")
	return
}

func isNilPointerInsideInterfaceValue(value interface{}) bool {
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Ptr:
		return v.IsNil()
	}
	return false
}


// The actual value must be <true>.
func IsTrue(actual interface{}, _ interface{}) (match bool, pos Message, neg Message, err os.Error) {
	match = actual.(bool) == true
	pos = Messagef(actual, "is <true>")
	neg = Messagef(actual, "is NOT <true>")
	return
}


// The actual value must be <false>.
func IsFalse(actual interface{}, _ interface{}) (match bool, pos Message, neg Message, err os.Error) {
	match = actual.(bool) == false
	pos = Messagef(actual, "is <false>")
	neg = Messagef(actual, "is NOT <false>")
	return
}


// The actual value must satisfy the given criteria.
func Satisfies(actual interface{}, criteria interface{}) (match bool, pos Message, neg Message, err os.Error) {
	match = criteria.(bool) == true
	pos = Messagef(actual, "satisfies the criteria")
	neg = Messagef(actual, "does NOT satisfy the criteria")
	return
}


// The actual value must satisfy the given criteria as described
// by the given description.
//
// Like Satisfies but allows to name the criteria. The description
// and its args are formatted using fmt.Sprintf().
func Is(description string, args ...interface{}) Matcher {
	return func(actual interface{}, criteria interface{}) (match bool, pos Message, neg Message, err os.Error) {
		match = criteria.(bool) == true
		pos = Messagef(actual, "is %v", fmt.Sprintf(description, args...))
		neg = Messagef(actual, "is NOT %v", fmt.Sprintf(description, args...))
		return
	}
}


// The actual value must be within delta from the expected value.
func IsWithin(delta float64) Matcher {
	return func(actual_ interface{}, expected_ interface{}) (match bool, pos Message, neg Message, err os.Error) {
		actual, err := toFloat64(actual_)
		if err != nil {
			return
		}
		expected, err := toFloat64(expected_)
		if err != nil {
			return
		}

		match = math.Fabs(expected-actual) < delta
		pos = Messagef(actual, "is within %v ± %v", expected, delta)
		neg = Messagef(actual, "is NOT within %v ± %v", expected, delta)
		return
	}
}

func toFloat64(actual interface{}) (result float64, err os.Error) {
	switch v := actual.(type) {
	case float32:
		result = float64(v)
	case float64:
		result = float64(v)
	default:
		err = Errorf("type error: expected a float, but was “%v” of type “%T”", actual, actual)
	}
	return
}


// The actual collection must contain the expected value.
func Contains(actual_ interface{}, expected interface{}) (match bool, pos Message, neg Message, err os.Error) {
	actual, err := toArray(actual_)
	if err != nil {
		return
	}

	match = arrayContains(actual, expected)
	pos = Messagef(actual, "contains “%v”", expected)
	neg = Messagef(actual, "does NOT contain “%v”", expected)
	return
}

func toArray(values interface{}) ([]interface{}, os.Error) {

	// vector to array
	if vector, ok := values.(*vector.Vector); ok {
		return toArray(*vector)
	}

	result := new(vector.Vector)

	// list to array
	if list, ok := values.(*list.List); ok {
		for e := list.Front(); e != nil; e = e.Next() {
			result.Push(e.Value)
		}
		return *result, nil
	}

	switch v := reflect.ValueOf(values); v.Kind() {

	// array to array (copy)
	case reflect.Array, reflect.Slice:
		arr := v
		for i := 0; i < arr.Len(); i++ {
			obj := arr.Index(i).Interface()
			result.Push(obj)
		}

	// channel to array
	case reflect.Chan:
		ch := v
		for {
			if x, ok := ch.Recv(); ok {
				obj := x.Interface()
				result.Push(obj)
			} else {
				break
			}
		}

	// unknown type
	default:
		return nil, Errorf("type error: expected a collection type, but was “%v” of type “%T”", values, values)
	}
	return *result, nil
}

func arrayContains(haystack []interface{}, needle interface{}) bool {
	_, found := findIndex(haystack, needle)
	return found
}

func findIndex(haystack []interface{}, needle interface{}) (idx int, found bool) {
	for i := 0; i < len(haystack); i++ {
		if areEqual(haystack[i], needle) {
			return i, true
		}
	}
	return -1, false
}


// The actual collection must contain all expected elements,
// but it may contain also other non-expected elements.
// The order of elements is not significant.
func ContainsAll(actual_ interface{}, expected_ interface{}) (match bool, pos Message, neg Message, err os.Error) {
	actual, err := toArray(actual_)
	if err != nil {
		return
	}
	expected, err := toArray(expected_)
	if err != nil {
		return
	}

	containsAll := true
	for i := 0; i < len(expected); i++ {
		if !arrayContains(actual, expected[i]) {
			containsAll = false
			break
		}
	}

	match = containsAll
	pos = Messagef(actual, "contains all of “%v”", expected)
	neg = Messagef(actual, "does NOT contain all of “%v”", expected)
	return
}


// The actual collection must contain at least one of the expected elements.
func ContainsAny(actual_ interface{}, expected_ interface{}) (match bool, pos Message, neg Message, err os.Error) {
	actual, err := toArray(actual_)
	if err != nil {
		return
	}
	expected, err := toArray(expected_)
	if err != nil {
		return
	}

	containsAny := false
	for i := 0; i < len(expected); i++ {
		if arrayContains(actual, expected[i]) {
			containsAny = true
			break
		}
	}

	match = containsAny
	pos = Messagef(actual, "contains any of “%v”", expected)
	neg = Messagef(actual, "does NOT contain any of “%v”", expected)
	return
}


// The actual collection must contain all expected elements and nothing else.
// The order of elements is not significant.
func ContainsExactly(actual_ interface{}, expected_ interface{}) (match bool, pos Message, neg Message, err os.Error) {
	actual, err := toArray(actual_)
	if err != nil {
		return
	}
	expected, err := toArray(expected_)
	if err != nil {
		return
	}

	containsAll := true
	remaining := new(vector.Vector)
	remaining.AppendVector((*vector.Vector)(&actual))
	for i := 0; i < len(expected); i++ {
		if idx, found := findIndex(*remaining, expected[i]); found {
			remaining.Delete(idx)
		} else {
			containsAll = false
			break
		}
	}

	match = containsAll && remaining.Len() == 0
	pos = Messagef(actual, "contains exactly “%v”", expected)
	neg = Messagef(actual, "does NOT contain exactly “%v”", expected)
	return
}


// The actual collection must contain all expected elements, in the same order, and nothing else.
func ContainsInOrder(actual_ interface{}, expected_ interface{}) (match bool, pos Message, neg Message, err os.Error) {
	actual, err := toArray(actual_)
	if err != nil {
		return
	}
	expected, err := toArray(expected_)
	if err != nil {
		return
	}

	containsInOrder := len(actual) == len(expected)
	for i := 0; i < len(actual) && i < len(expected); i++ {
		if !areEqual(actual[i], expected[i]) {
			containsInOrder = false
		}
	}

	match = containsInOrder
	pos = Messagef(actual, "contains in order “%v”", expected)
	neg = Messagef(actual, "does NOT contain in order “%v”", expected)
	return
}


// The actual collection must contain all expected objects, in the same order,
// but it may contain also other non-expected objects.
// For example [1, 2, 2, 3, 4] contains in partial order [1, 2, 3].
// See http://en.wikipedia.org/wiki/Partial_order for further information.
func ContainsInPartialOrder(actual_ interface{}, expected_ interface{}) (match bool, pos Message, neg Message, err os.Error) {
	actual, err := toArray(actual_)
	if err != nil {
		return
	}
	expected, err := toArray(expected_)
	if err != nil {
		return
	}

	containsInPartialOrder := true
	for ie, ia := 0, 0; ie < len(expected); {
		if ia >= len(actual) {
			containsInPartialOrder = false
			break
		}
		if areEqual(actual[ia], expected[ie]) {
			ie++
			ia++
		} else {
			ia++
		}
	}

	match = containsInPartialOrder
	pos = Messagef(actual, "contains in partial order “%v”", expected)
	neg = Messagef(actual, "does NOT contain in partial order “%v”", expected)
	return
}


// Special value thrown to indicate that there was no error
var noError interface{} = &noErrorT{}

type noErrorT struct{}

func (noErr noErrorT) String() string {
	return "<gospec: not an error condition>"
}

// Runs code() and returns whatever value it threw by calling panic(value).
// Or if panic(value) was not called, returns a special value indicating
// that no error happened, which can be tested with IsNoError.
func Catch(code func()) (result interface{}) {
	defer func() {
		result = recover()
	}()
	code()
	panic(noError)
}

// Checks whether a call to Catch() returned a special value indicating
// that panic() was not called by the code. Used together with either
// Panic or Catch, but not both at the same time. 
//
// May be useful in cases when Panic and RunsNormally are not sufficient
// and you need to do more complicated things with the error value.
func IsNoError(actual interface{}, _ interface{}) (match bool, pos Message, neg Message, err os.Error) {
	match = (actual == noError)
	pos = Messagef(actual, "is not an error")
	neg = Messagef(actual, "IS an error")
	return
}


// The actual func() must panic with a value which matches the expected value
// using the provided matcher.
func Panic(matcher Matcher) Matcher {
	return func(actual interface{}, expected interface{}) (match bool, pos Message, neg Message, err os.Error) {
		closure, ok := actual.(func())
		if !ok {
			err = Errorf("type error: expected a func(), but was “%v” of type “%T”", actual, actual)
			return
		}

		thrownValue := Catch(closure)
		match, matcherPos, _, err := matcher(thrownValue, expected)
		pos = Messagef(matcherPos.Actual(), fmt.Sprintf("panics with a value that %v", matcherPos.Expectation()))
		neg = Messagef(matcherPos.Actual(), fmt.Sprintf("does NOT panic with a value that %v", matcherPos.Expectation()))
		return
	}
}


// The actual func() must not panic.
//
// This is the same as Panic(IsNoError) except that it prints nicer messages.
func RunsNormally(actual interface{}, expected interface{}) (match bool, pos Message, neg Message, err os.Error) {
	match, _, _, err = Panic(IsNoError)(actual, expected)
	pos = Messagef(actual, "runs normally (i.e. does NOT panic)")
	neg = Messagef(actual, "does NOT run normally (i.e. panics)")
	return
}
