// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"reflect"
	"runtime"
)


func functionName(function interface{}) string {
	fval := reflect.NewValue(function).(*reflect.FuncValue)
	f := runtime.FuncForPC(fval.Get())
	if f != nil {
		return f.Name()
	}
	return "<unknown function>"
}
