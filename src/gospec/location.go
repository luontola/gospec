// Copyright © 2009 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt"
	filepath "path"
	"runtime"
)


type Location struct {
	File string
	Line int
}

func currentLocation() *Location {
	return newLocation(1)
}

func callerLocation() *Location {
	return newLocation(2)
}

func newLocation(n int) *Location {
	if _, file, line, ok := runtime.Caller(n + 1); ok {
		return &Location{filename(file), line}
	}
	return nil
}

func filename(path string) string {
	_, file := filepath.Split(path)
	return file
}

func (loc *Location) String() string {
	if loc == nil {
		return "Unknown File"
	}
	return fmt.Sprintf("%v:%v", loc.File, loc.Line)
}

