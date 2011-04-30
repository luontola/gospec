// Copyright © 2009-2011 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"reflect"
	"runtime"
)

const unknownFunction = "<unknown function>"

func functionName(function interface{}) string {
	if f := functionToFunc(function); f != nil {
		return f.Name()
	}
	return unknownFunction
}

func functionToFunc(function interface{}) *runtime.Func {
	fval := reflect.ValueOf(function)
	return runtime.FuncForPC(fval.Pointer())
}
