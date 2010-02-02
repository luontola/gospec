// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt"
	"os"
)


type NewMatcher func(actual interface{}, expected interface{}) (ok bool, pos os.Error, neg os.Error)

func Not(matcher NewMatcher) NewMatcher {
	return func(actual interface{}, expected interface{}) (ok bool, pos os.Error, neg os.Error) {
		ok, pos, neg = matcher(actual, expected)
		ok = !ok
		pos, neg = neg, pos
		return
	}
}


func Equals(actual interface{}, expected interface{}) (ok bool, pos os.Error, neg os.Error) {
	ok = areEqual(actual, expected)
	// TODO: change the messages to following?
	// '%v' should equal '%v', but it did not
	// '%v' should NOT equal '%v', but it did
	pos = Errorf("Expected '%v' but was '%v'", expected, actual)
	neg = Errorf("Did not expect '%v' but was '%v'", expected, actual)
	return
}

func areEqual(a interface{}, b interface{}) bool {
	if a2, ok := a.(Equality); ok {
		return a2.Equals(b)
	}
	// TODO: pointer equality
	return a == b
}

type Equality interface {
	Equals(other interface{}) bool
}


// Helpers

func Errorf(format string, args ...) os.Error {
	return lazyString(func() string {
		return fmt.Sprintf(format, args)
	})
}

type lazyString func() string

func (this lazyString) String() string {
	return this()
}

