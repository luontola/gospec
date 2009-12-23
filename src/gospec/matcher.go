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
	required  bool
	Should    *Matcher
	ShouldNot *Matcher
	Must      *Matcher
	MustNot   *Matcher
}

func newMatcher(actual interface{}, log errorLogger) *Matcher {
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
		m.log       = log
		m.Should    = posOpt
		m.ShouldNot = negOpt
		m.Must      = posReq
		m.MustNot   = negReq
	}
	return posOpt
}

func (m *Matcher) Equal(expected interface{}) {
	if m.fails(areEqual(expected, m.actual)) {
		m.addError(expected, m.actual)
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

func (m *Matcher) addError(expected interface{}, actual interface{}) {
	var message string
	if m.negation {
		message = fmt.Sprintf("Did not expect '%v' but was '%v'", expected, actual)
	} else {
		message = fmt.Sprintf("Expected '%v' but was '%v'", expected, actual)
	}
	if m.required {
		m.log.AddFatalError(message)
	} else {
		m.log.AddError(message)
	}
}


// Helpers

type Equality interface {
	Equals(other interface{}) bool
}

type errorLogger interface {
	AddError(message string)
	AddFatalError(message string)
}

