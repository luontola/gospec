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
	StackTrace []*runtime.Func
}

func (this *exception) String() string {
	return fmt.Sprintf("panic: %v", this.Cause)
}


func recoverOnPanic(f func()) (err *exception) {
	defer func() {
		if cause := recover(); cause != nil {
			st := stackTraceOfPanic()
			st = cutStackTraceAt(functionToFunc(recoverOnPanic), st)
			err = &exception{cause, st}
		}
	}()
	f()
	return
}

func stackTraceOfPanic() []*runtime.Func {
	callers := make([]uintptr, 16)
	for {
		count := runtime.Callers(4, callers)
		if count == len(callers) {
			callers = make([]uintptr, len(callers)*2)
		} else {
			callers = callers[0:count]
			break
		}
	}
	return asFuncArray(callers)
}

func asFuncArray(ptrs []uintptr) []*runtime.Func {
	result := make([]*runtime.Func, len(ptrs))
	for i, ptr := range ptrs {
		result[i] = runtime.FuncForPC(ptr)
	}
	return result
}

func cutStackTraceAt(cutpoint *runtime.Func, stacktrace []*runtime.Func) []*runtime.Func {
	for i, f := range stacktrace {
		if f.Entry() == cutpoint.Entry() {
			return stacktrace[0:i]
		}
	}
	return stacktrace
}
