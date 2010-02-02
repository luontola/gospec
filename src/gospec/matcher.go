// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
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


// Matchers

// The actual value must equal the expected value. For primitives the equality
// operator is used. All other objects must implement the Equality interface.
func (m *Matcher) Equal(expected interface{}) {
	m.compareUsing(Equals, expected)
}

// The actual value must satisfy the given criteria.
func (m *Matcher) Be(criteria bool) {
	m.compareUsing(Satisfies, criteria)
}

// The actual value must be within delta from the expected value.
func (m *Matcher) BeNear(expected float64, delta float64) {
	m.compareUsing(IsWithin(delta), expected)
}

// The actual collection (array, slice, iterator/channel) must contain the expected value.
func (m *Matcher) Contain(expected interface{}) {
	m.compareUsing(Contains, expected)
}

