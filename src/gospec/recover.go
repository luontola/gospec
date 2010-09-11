// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"fmt"
	"runtime"
)


type exception struct {
	Cause      interface{}
	StackTrace []*Location
}

func (this *exception) ToError() *Error {
	return newError(OtherError, this.String(), "", this.StackTrace)
}

func (this *exception) String() string {
	return fmt.Sprintf("panic: %v", this.Cause)
}


func recoverOnPanic(f func()) (err *exception) {
	defer func() {
		if cause := recover(); cause != nil {
			callers := stackTraceOfPanic()
			callers = cutStackTraceAt(recoverOnPanic, callers)
			err = &exception{cause, asLocationArray(callers)}
		}
	}()
	f()
	return
}

func stackTraceOfPanic() []uintptr {
	// When changing this method, remember to test the array resizing code
	// by temporarily setting the initial array size to 1.
	callers := make([]uintptr, 16)
	for {
		// Magic number for correct operation when called from recoverOnPanic()
		count := runtime.Callers(4, callers)
		if count == len(callers) {
			callers = make([]uintptr, len(callers)*2)
		} else {
			callers = callers[0:count]
			break
		}
	}
	return callers
}

func cutStackTraceAt(cutpoint_ interface{}, callers []uintptr) []uintptr {
	cutpoint := functionToFunc(cutpoint_).Entry()

	for i, ptr := range callers {
		current := runtime.FuncForPC(ptr).Entry()
		if current == cutpoint {
			return callers[0:i]
		}
	}
	return callers
}

func asLocationArray(pcs []uintptr) []*Location {
	result := make([]*Location, len(pcs))
	for i, pc := range pcs {
		result[i] = locationForPC(pcOfWhereCallWasMade(pc))
	}
	return result
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
