// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
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
	if pc, _, _, ok := runtime.Caller(n + 1); ok {
		return locationForPC(pc)
	}
	return nil
}

func locationForPC(pc uintptr) *Location {
	pc = pcOfWhereCallWasMade(pc)
	f := runtime.FuncForPC(pc)
	name := f.Name()
	file, line := f.FileLine(pc)
	return &Location{name, file, line}
}

// Quoted from http://code.google.com/p/go/issues/detail?id=1100
//   "It's a subtle thing, but runtime.Callers returns the return PCs
//   going up the stack.  The return PCs are the PCs of the instruction
//   that the call returns to, not the call itself.  The code that formats
//   the traceback for a crash subtracts 1 from each PC before translating
//   it to a line number.  The equivalent change in your code would be to
//   call f.FileLine(pc - 1)."
func pcOfWhereCallWasMade(pcOfWhereCallReturnsTo uintptr) uintptr {
	return pcOfWhereCallReturnsTo - 1
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
