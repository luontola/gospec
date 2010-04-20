// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"reflect"
	"runtime"
)


func functionName(function interface{}) string {
	f := functionToFunc(function)
	if f != nil {
		return f.Name()
	}
	return "<unknown function>"
}

func functionToFunc(function interface{}) *runtime.Func {
	fval := reflect.NewValue(function).(*reflect.FuncValue)
	return runtime.FuncForPC(fval.Get())
}
