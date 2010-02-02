// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"container/list"
	"fmt"
	"math"
	"reflect"
)


type Matcher struct {
	actual    interface{}
	location  *Location
	log       errorLogger
	negation  bool
	required  bool
	Should    *Matcher
	ShouldNot *Matcher
	Must      *Matcher
	MustNot   *Matcher
}

func newMatcher(actual interface{}, location *Location, log errorLogger) *Matcher {
	posOpt := new(Matcher)
	posOpt.negation = false
	posOpt.required = false
	
	posReq := new(Matcher)
	posReq.negation = false
	posReq.required = true
	
	negOpt := new(Matcher)
	negOpt.negation = true
	negOpt.required = false
	
	negReq := new(Matcher)
	negReq.negation = true
	negReq.required = true
	
	all := [...]*Matcher{posOpt, posReq, negOpt, negReq}
	for _, m := range all {
		m.actual    = actual
		m.location  = location
		m.log       = log
		m.Should    = posOpt
		m.ShouldNot = negOpt
		m.Must      = posReq
		m.MustNot   = negReq
	}
	return posOpt
}

func (m *Matcher) compareUsing(matcher NewMatcher, expected interface{}) {
	if m.negation {
		matcher = Not(matcher)
	}
	if ok, error, _ := matcher(m.actual, expected); !ok {
		m.addError(error.String())
	}
}

func (m *Matcher) addError(message string) {
	error := newError(message, m.location)
	if m.required {
		m.log.AddFatalError(error)
	} else {
		m.log.AddError(error)
	}
}


// The actual value must equal the expected value. For primitives the equality
// operator is used. All other objects must implement the Equality interface.
func (m *Matcher) Equal(expected interface{}) {
	m.compareUsing(Equals, expected)
}

// The actual value must satisfy the given criteria.
func (m *Matcher) Be(criteria bool) {
	if m.fails(criteria) {
		m.addError(fmt.Sprintf("Criteria not satisfied by '%v'", m.actual))
	}
}

// The actual value must be within delta from the expected value.
func (m *Matcher) BeNear(expected float64, delta float64) {
	var actual float64
	switch v := m.actual.(type) {
	case float:
		actual = float64(v)
	case float32:
		actual = float64(v)
	case float64:
		actual = float64(v)
	default:
		m.addError(fmt.Sprintf("Expected a float, but was '%v' of type '%T'", m.actual, m.actual))
		return
	}
	if m.fails(math.Fabs(expected - actual) < delta) {
		if m.negation {
			m.addError(fmt.Sprintf("Did not expect '%v' ± %v but was '%v'", expected, delta, actual))
		} else {
			m.addError(fmt.Sprintf("Expected '%v' ± %v but was '%v'", expected, delta, actual))
		}
	}
}

// The actual collection (array, slice, iterator/channel) must contain the expected value.
func (m *Matcher) Contain(expected interface{}) {
	switch v := reflect.NewValue(m.actual).(type) {
	
	case reflect.ArrayOrSliceValue:
		arr := v
		contains := false
		for i := 0; i < arr.Len(); i++ {
			other :=  arr.Elem(i).Interface()
			if areEqual(expected, other) {
				contains = true
				break
			}
		}
		// TODO: remove duplication (maybe better done after there are more containment matchers)
		if m.fails(contains) {
			if m.negation {
				m.addError(fmt.Sprintf("Did not expect '%v' to be in '%v' but it was", expected, m.actual))
			} else {
				m.addError(fmt.Sprintf("Expected '%v' to be in '%v' but it was not", expected, m.actual))
			}
		}
	
	case *reflect.ChanValue:
		ch := v
		contains := false
		list := list.New()
		for {
			other := ch.Recv().Interface()
			if ch.Closed() {
				break
			}
			if areEqual(expected, other) {
				contains = true
			}
			list.PushBack(other)
		}
		if m.fails(contains) {
			actual := listToArray(list)
			if m.negation {
				m.addError(fmt.Sprintf("Did not expect '%v' to be in '%v' but it was", expected, actual))
			} else {
				m.addError(fmt.Sprintf("Expected '%v' to be in '%v' but it was not", expected, actual))
			}
		}
	
	default:
		m.addError(fmt.Sprintf("Unknown type '%T', not iterable: %v", m.actual, m.actual))
	}
}

// TODO: ContainAll - The actual collection must contain all given elements. The order of elements is not significant.
// TODO: ContainAny - The actual collection must contain at least one element from the given collection.
// TODO: ContainExactly - The actual collection must contain exactly the same elements as in the given collection. The order of elements is not significant.
// TODO: ContainInOrder - The actual collection must contain exactly the same elements as in the given collection, and they must be in the same order.
// TODO: ContainInPartialOrder - The actual collection can hold other objects, but the objects which are common in both collections must be in the same order. The actual collection can also repeat some elements. For example [1, 2, 2, 3, 4] contains in partial order [1, 2, 3]. See Wikipedia <http://en.wikipedia.org/wiki/Partial_order> for further information.


func listToArray(list *list.List) []interface{} {
	arr := make([]interface{}, list.Len())
	i := 0
	for e := list.Front(); e != nil; e = e.Next() {
		arr[i] = e.Value
		i++
	}
	return arr
}

func (m *Matcher) fails(ok bool) bool { // TODO: remove
	return m.negation == ok
}

