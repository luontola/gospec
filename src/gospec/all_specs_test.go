// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gospec

import (
	"nanospec"
	"testing"
)


func TestAllSpecs(t *testing.T) {
	nanospec.Run(t, ConcurrencySpec)
	nanospec.Run(t, ExecutingSpecsSpec)
	nanospec.Run(t, ExecutionModelSpec)
	nanospec.Run(t, ExpectationsSpec)
	nanospec.Run(t, FuncNameSpec)
	nanospec.Run(t, LocationSpec)
	nanospec.Run(t, MatcherMessagesSpec)
	nanospec.Run(t, MatchersSpec)
	nanospec.Run(t, PrinterSpec)
	nanospec.Run(t, ResultsSpec)
}
