// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import ()


type Error struct {
	Message    string
	StackTrace []*Location
}

func newError(message string, stacktrace []*Location) *Error {
	return &Error{message, stacktrace}
}

func (this *Error) equals(that *Error) bool {
	return this.Message == that.Message &&
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
	AddFatalError(error *Error)
}
