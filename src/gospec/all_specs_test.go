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
}
