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

func (this *Location) equals(that *Location) bool {
	return this.File == that.File &&
	       this.Line == that.Line
}

func (this *Location) String() string {
	if this == nil {
		return "Unknown File"
	}
	return fmt.Sprintf("%v:%v", this.File, this.Line)
}

