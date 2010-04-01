// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"testing"
)


func Test__Get_the_location_of_current_method(t *testing.T) {
	loc := currentLocation()
	assertHasPrefix("location_test.go:", loc.String(), t)
}

func Test__Get_the_location_of_calling_method(t *testing.T) {
	loc := callerLocation()
	assertHasPrefix("testing.go:", loc.String(), t)
}

func Test__Failing_to_get_the_location_fails_grafecully_with_an_error_message(t *testing.T) {
	loc := newLocation(1000)
	assertEquals("Unknown File", loc.String(), t)
}

func Test__Calls_to_newLocation_are_synced_with_the_helper_methods(t *testing.T) {
	assertEquals(currentLocation().String(), newLocation(0).String(), t)
	assertEquals(callerLocation().String(), newLocation(1).String(), t)
}
