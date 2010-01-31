// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
)


type Error struct {
	Message  string
	Location *Location
}

func newError(message string, location *Location) *Error {
	return &Error{message, location}
}

func (this *Error) equals(that *Error) bool {
	return this.Message == that.Message &&
	       this.Location.equals(that.Location)
}

func (this *Error) String() string {
	return this.Location.String() + " " + this.Message
}

type errorLogger interface {
	AddError(error *Error)
	AddFatalError(error *Error)
}

