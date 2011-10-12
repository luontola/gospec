// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import ()

type ErrorType int

const (
	ExpectFailed ErrorType = iota
	AssumeFailed
	OtherError
)

type Error struct {
	Type       ErrorType
	Message    string
	Actual     string
	StackTrace []*Location
}

func newError(errortype ErrorType, message string, actual string, stacktrace []*Location) *Error {
	return &Error{errortype, message, actual, stacktrace}
}

func (this *Error) equals(that *Error) bool {
	return this.Message == that.Message &&
		this.Actual == that.Actual &&
		stackTracesEqual(this.StackTrace, that.StackTrace)
}

func stackTracesEqual(a []*Location, b []*Location) bool {
	if len(a) != len(b) {
		return false
	}
	for i, _ := range a {
		if !a[i].equals(b[i]) {
			return false
		}
	}
	return true
}

func (this *Error) String() string {
	s := this.Message
	for _, loc := range this.StackTrace {
		s += "\nat " + loc.String()
	}
	return s
}

type errorLogger interface {
	AddError(error *Error)
}

type ratedErrorLogger interface {
	errorLogger
	AddFatalError(error *Error)
}
