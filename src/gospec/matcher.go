// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
)


// TODO: remove when the c.Then() syntax is removed

type MatcherBuilder struct {
	actual    interface{}
	location  *Location
	log       errorLogger
	negation  bool
	required  bool
	Should    *MatcherBuilder
	ShouldNot *MatcherBuilder
	Must      *MatcherBuilder
	MustNot   *MatcherBuilder
}

func newMatcherBuilder(actual interface{}, location *Location, log errorLogger) *MatcherBuilder {
	posOpt := new(MatcherBuilder)
	posOpt.negation = false
	posOpt.required = false
	
	posReq := new(MatcherBuilder)
	posReq.negation = false
	posReq.required = true
	
	negOpt := new(MatcherBuilder)
	negOpt.negation = true
	negOpt.required = false
	
	negReq := new(MatcherBuilder)
	negReq.negation = true
	negReq.required = true
	
	all := [...]*MatcherBuilder{posOpt, posReq, negOpt, negReq}
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

func (m *MatcherBuilder) compareUsing(matcher Matcher, expected interface{}) {
	if m.negation {
		matcher = Not(matcher)
	}
	ok, pos, _, err := matcher(m.actual, expected)
	if err != nil {
		m.addError(err.String())
	} else if !ok {
		m.addError(pos.String())
	}
}

func (m *MatcherBuilder) addError(message string) {
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
func (m *MatcherBuilder) Equal(expected interface{}) {
	m.compareUsing(Equals, expected)
}

// The actual value must satisfy the given criteria.
func (m *MatcherBuilder) Be(criteria bool) {
	m.compareUsing(Satisfies, criteria)
}

// The actual value must be within delta from the expected value.
func (m *MatcherBuilder) BeNear(expected float64, delta float64) {
	m.compareUsing(IsWithin(delta), expected)
}

// The actual collection (array, slice, iterator/channel) must contain the expected value.
func (m *MatcherBuilder) Contain(expected interface{}) {
	m.compareUsing(Contains, expected)
}

