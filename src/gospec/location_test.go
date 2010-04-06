// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"nanospec"
	"strings"
)


func LocationSpec(c nanospec.Context) {

	c.Specify("Location of the current method can be found", func() {
		loc := currentLocation()
		c.Expect(loc).Satisfies(strings.HasPrefix(loc.String(), "location_test.go:"))
	})
	c.Specify("Location of the calling method can be found", func() {
		loc := callerLocation()
		c.Expect(loc).Satisfies(strings.HasPrefix(loc.String(), "context.go:"))
	})
	c.Specify("When failing to get the location, it will fail gracefully with an error message", func() {
		loc := newLocation(1000)
		c.Expect(loc.String()).Equals("Unknown File")
	})
	c.Specify("Calls to newLocation are synced with the helper methods", func() {
		c.Expect(newLocation(0).String()).Equals(currentLocation().String())
		c.Expect(newLocation(1).String()).Equals(callerLocation().String())
	})
}
