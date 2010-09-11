// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"nanospec"
	"runtime"
)


func LocationSpec(c nanospec.Context) {

	c.Specify("Location of the current method can be found", func() {
		loc := currentLocation() // line 16
		c.Expect(loc.FileName()).Equals("location_test.go")
		c.Expect(loc.Line()).Equals(16)
	})
	c.Specify("Location of the calling method can be found", func() {
		loc := callerLocation()
		c.Expect(loc.FileName()).Equals("context.go")
	})
	c.Specify("The name of the method is provided", func() {
		loc := methodWhereLocationIsCalled()
		c.Expect(loc.Name()).Equals("gospec.methodWhereLocationIsCalled")
	})
	c.Specify("Calls to newLocation are synced with the helper methods", func() {
		c.Expect(newLocation(0).Name()).Equals(currentLocation().Name())
		c.Expect(newLocation(0).File()).Equals(currentLocation().File())
		c.Expect(newLocation(1).Name()).Equals(callerLocation().Name())
		c.Expect(newLocation(1).File()).Equals(callerLocation().File())
	})
	c.Specify("Program Counters can be converted to Locations", func() {
		expectedLine := currentLocation().Line() + 1
		pc, _, _, _ := runtime.Caller(0)
		loc := locationForPC(pc)
		c.Expect(loc.FileName()).Equals("location_test.go")
		c.Expect(loc.Line()).Equals(expectedLine)
	})
}

func methodWhereLocationIsCalled() *Location {
	return currentLocation()
}
