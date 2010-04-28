// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt"
	filepath "path"
	"runtime"
)


type Location struct {
	name string
	file string
	line int
}

func currentLocation() *Location {
	return newLocation(1)
}

func callerLocation() *Location {
	return newLocation(2)
}

func newLocation(n int) *Location {
	if pc, file, line, ok := runtime.Caller(n + 1); ok {
		name := functionNameFromPC(pc)
		return &Location{name, file, line}
		// TODO: replace with the following code when the bug in runtime.Func.FileLine() is fixed
		// return locationForPC(pc)
	}
	return nil
}

func locationForPC(pc uintptr) *Location {
	f := runtime.FuncForPC(pc)
	name := f.Name()
	file, line := f.FileLine(pc)
	return &Location{name, file, line}
}

func (this *Location) Name() string     { return this.name }
func (this *Location) File() string     { return this.file }
func (this *Location) FileName() string { return filename(this.file) }
func (this *Location) Line() int        { return this.line }

func filename(path string) string {
	_, file := filepath.Split(path)
	return file
}

func (this *Location) equals(that *Location) bool {
	return this.name == that.name &&
		this.file == that.file &&
		this.line == that.line
}

func (this *Location) String() string {
	return fmt.Sprintf("%v:%v", this.FileName(), this.Line())
}
