// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt"
)


type Matcher struct {
	actual    interface{}
	log       errorLogger
	negation  bool
	Should    *Matcher
	ShouldNot *Matcher
}

func newMatcher(actual interface{}, log errorLogger) *Matcher {
	var common Matcher
	common.actual = actual
	common.log = log
	
	pos := common
	pos.negation = false
	
	neg := common
	neg.negation = true
	
	pos.Should = &pos
	pos.ShouldNot = &neg
	neg.Should = &pos
	neg.ShouldNot = &neg
	
	return &pos
}

func (m *Matcher) Equal(expected interface{}) {
	if m.fails(areEqual(expected, m.actual)) {
		m.logError(expected, m.actual)
	}
}

func areEqual(a interface{}, b interface{}) bool {
	if a2, ok := a.(Equality); ok {
		return a2.Equals(b)
	}
	return a == b
}

func (m *Matcher) fails(ok bool) bool {
	return m.negation == ok
}

func (m *Matcher) logError(expected interface{}, actual interface{}) {
	if !m.negation {
		m.log.logError(fmt.Sprintf("Expected '%v' but was '%v'", expected, actual))
	}
	if m.negation {
		m.log.logError(fmt.Sprintf("Did not expect '%v' but was '%v'", expected, actual))
	}
}


// Helpers

type Equality interface {
	Equals(other interface{}) bool
}

type errorLogger interface {
	logError(message string)
}

