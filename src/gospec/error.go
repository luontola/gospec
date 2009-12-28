// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
)


type Error struct {
	Message string
}

func newError(message string) *Error {
	return &Error{message}
}


type errorLogger interface {
	AddError(error *Error)
	AddFatalError(error *Error)
}

